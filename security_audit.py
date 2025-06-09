#!/usr/bin/env python3
"""
Security audit script for MoneyMaker application
Performs comprehensive security checks and generates a report
"""

import os
import sys
import json
import subprocess
import re
from pathlib import Path

class SecurityAuditor:
    """Security audit tool for the MoneyMaker application"""
    
    def __init__(self, project_root):
        self.project_root = Path(project_root)
        self.issues = []
        self.recommendations = []
        self.passed_checks = []
        
    def log_issue(self, severity, category, description, file_path=None):
        """Log a security issue"""
        self.issues.append({
            'severity': severity,
            'category': category,
            'description': description,
            'file': file_path
        })
        
    def log_recommendation(self, category, description):
        """Log a security recommendation"""
        self.recommendations.append({
            'category': category,
            'description': description
        })
        
    def log_passed(self, check_name):
        """Log a passed security check"""
        self.passed_checks.append(check_name)
        
    def check_sensitive_files(self):
        """Check for sensitive files that shouldn't be committed"""
        sensitive_patterns = [
            '.env',
            '*.key',
            '*.pem',
            '*.p12',
            '*.pfx',
            'secrets.txt',
            'passwords.txt',
            'config.ini',
            'database.db',
            '*.sqlite',
            '*.log'
        ]
        
        # Patterns to exclude (legitimate files)
        exclude_patterns = [
            '**/certifi/cacert.pem',  # SSL certificates from certifi package
            '**/pip/_vendor/certifi/cacert.pem',  # SSL certificates from pip's vendored certifi
            '**/.venv/**',  # Virtual environment files
            '**/venv/**',   # Virtual environment files
        ]
        
        for pattern in sensitive_patterns:
            files = list(self.project_root.glob(f"**/{pattern}"))
            for file in files:
                if file.name != '.env.example':  # Allow example files
                    # Check if file should be excluded
                    should_exclude = False
                    for exclude_pattern in exclude_patterns:
                        if file.match(exclude_pattern) or str(file).find('site-packages') != -1:
                            should_exclude = True
                            break
                    
                    if not should_exclude:
                        self.log_issue(
                            'HIGH', 
                            'Sensitive Files', 
                            f'Potentially sensitive file found: {file.relative_to(self.project_root)}',
                            str(file)
                        )
                    
        if not any(issue['category'] == 'Sensitive Files' for issue in self.issues):
            self.log_passed('Sensitive files check')
            
    def check_gitignore(self):
        """Check .gitignore for proper exclusions"""
        gitignore_path = self.project_root / '.gitignore'
        
        if not gitignore_path.exists():
            self.log_issue('HIGH', 'Git Security', '.gitignore file missing')
            return
            
        with open(gitignore_path, 'r') as f:
            gitignore_content = f.read()
            
        required_patterns = [
            '.env',
            '__pycache__',
            '*.pyc',
            '*.log',
            'venv/',
            '.venv/',
            'node_modules/',
            '*.key',
            '*.pem',
            'secrets*'
        ]
        
        missing_patterns = []
        for pattern in required_patterns:
            if pattern not in gitignore_content:
                missing_patterns.append(pattern)
                
        if missing_patterns:
            self.log_issue(
                'MEDIUM',
                'Git Security',
                f'Missing .gitignore patterns: {", ".join(missing_patterns)}'
            )
        else:
            self.log_passed('Gitignore security patterns')
            
    def check_flask_security(self):
        """Check Flask security configuration"""
        app_py_path = self.project_root / 'src' / 'app.py'
        
        if not app_py_path.exists():
            self.log_issue('HIGH', 'Flask Security', 'app.py not found')
            return
            
        with open(app_py_path, 'r') as f:
            app_content = f.read()
            
        # Check for debug mode
        if 'debug=True' in app_content:
            self.log_issue('HIGH', 'Flask Security', 'Debug mode enabled in production')
        else:
            self.log_passed('Flask debug mode disabled')
            
        # Check for secret key
        if 'SECRET_KEY' in app_content:
            self.log_passed('Flask secret key configured')
        else:
            self.log_issue('HIGH', 'Flask Security', 'No SECRET_KEY configuration found')
            
        # Check for CSRF protection
        if 'CSRFProtect' in app_content:
            self.log_passed('CSRF protection enabled')
        else:
            self.log_issue('MEDIUM', 'Flask Security', 'CSRF protection not found')
            
        # Check for rate limiting
        if 'limiter' in app_content or 'Limiter' in app_content:
            self.log_passed('Rate limiting configured')
        else:
            self.log_issue('MEDIUM', 'Flask Security', 'Rate limiting not configured')
            
        # Check for security headers
        if 'X-Content-Type-Options' in app_content:
            self.log_passed('Security headers configured')
        else:
            self.log_issue('MEDIUM', 'Flask Security', 'Security headers not configured')
            
    def check_input_validation(self):
        """Check for input validation implementation"""
        app_py_path = self.project_root / 'src' / 'app.py'
        
        if not app_py_path.exists():
            return
            
        with open(app_py_path, 'r') as f:
            app_content = f.read()
            
        # Check for validation functions
        if 'validate_input' in app_content:
            self.log_passed('Input validation decorators found')
        else:
            self.log_issue('MEDIUM', 'Input Validation', 'Input validation decorators not found')
            
        # Check for regex validation
        if 're.match' in app_content:
            self.log_passed('Regex input validation implemented')
        else:
            self.log_issue('LOW', 'Input Validation', 'Regex validation not found')
            
    def check_logging_security(self):
        """Check logging security configuration"""
        logging_config_path = self.project_root / 'src' / 'utils' / 'logging_config.py'
        
        if not logging_config_path.exists():
            self.log_issue('MEDIUM', 'Logging Security', 'Logging configuration not found')
            return
            
        with open(logging_config_path, 'r') as f:
            logging_content = f.read()
            
        # Check for sensitive data filtering
        if 'SensitiveDataFilter' in logging_content:
            self.log_passed('Sensitive data filtering in logs')
        else:
            self.log_issue('MEDIUM', 'Logging Security', 'Sensitive data filtering not found')
            
        # Check for log rotation
        if 'RotatingFileHandler' in logging_content:
            self.log_passed('Log rotation configured')
        else:
            self.log_issue('LOW', 'Logging Security', 'Log rotation not configured')
            
    def check_dependencies(self):
        """Check dependencies for known vulnerabilities"""
        requirements_path = self.project_root / 'requirements.txt'
        
        if not requirements_path.exists():
            self.log_issue('MEDIUM', 'Dependencies', 'requirements.txt not found')
            return
            
        with open(requirements_path, 'r') as f:
            requirements_content = f.read()
            
        # Check for version pinning
        lines = requirements_content.strip().split('\n')
        unpinned_deps = []
        
        for line in lines:
            if line.strip() and not any(op in line for op in ['==', '>=', '<=', '>', '<', '~=']):
                unpinned_deps.append(line.strip())
                
        if unpinned_deps:
            self.log_issue(
                'MEDIUM',
                'Dependencies',
                f'Unpinned dependencies found: {", ".join(unpinned_deps)}'
            )
        else:
            self.log_passed('All dependencies version pinned')
            
        # Check for security-related packages
        security_packages = ['flask-limiter', 'flask-wtf', 'python-dotenv']
        missing_security_packages = []
        
        for package in security_packages:
            if package not in requirements_content:
                missing_security_packages.append(package)
                
        if missing_security_packages:
            self.log_issue(
                'LOW',
                'Dependencies',
                f'Recommended security packages missing: {", ".join(missing_security_packages)}'
            )
        else:
            self.log_passed('Security packages present')
            
    def check_environment_config(self):
        """Check environment configuration"""
        env_example_path = self.project_root / '.env.example'
        
        if not env_example_path.exists():
            self.log_issue('MEDIUM', 'Environment Config', '.env.example file missing')
            return
            
        with open(env_example_path, 'r') as f:
            env_content = f.read()
            
        # Check for important environment variables
        required_vars = ['SECRET_KEY', 'FLASK_DEBUG', 'FLASK_ENV']
        missing_vars = []
        
        for var in required_vars:
            if var not in env_content:
                missing_vars.append(var)
                
        if missing_vars:
            self.log_issue(
                'MEDIUM',
                'Environment Config',
                f'Missing environment variables in .env.example: {", ".join(missing_vars)}'
            )
        else:
            self.log_passed('Environment configuration complete')
            
        # Check for production defaults
        if 'FLASK_DEBUG=False' in env_content:
            self.log_passed('Production debug mode default')
        else:
            self.log_issue('HIGH', 'Environment Config', 'Debug mode not defaulted to False')
            
    def generate_report(self):
        """Generate security audit report"""
        report = {
            'summary': {
                'total_issues': len(self.issues),
                'high_severity': len([i for i in self.issues if i['severity'] == 'HIGH']),
                'medium_severity': len([i for i in self.issues if i['severity'] == 'MEDIUM']),
                'low_severity': len([i for i in self.issues if i['severity'] == 'LOW']),
                'passed_checks': len(self.passed_checks),
                'recommendations': len(self.recommendations)
            },
            'issues': self.issues,
            'passed_checks': self.passed_checks,
            'recommendations': self.recommendations
        }
        
        return report
        
    def run_audit(self):
        """Run complete security audit"""
        print("🔒 Running MoneyMaker Security Audit...")
        print("=" * 50)
        
        # Run all checks
        self.check_sensitive_files()
        self.check_gitignore()
        self.check_flask_security()
        self.check_input_validation()
        self.check_logging_security()
        self.check_dependencies()
        self.check_environment_config()
        
        # Generate and display report
        report = self.generate_report()
        
        print(f"\n📊 AUDIT SUMMARY")
        print(f"================")
        print(f"✅ Passed checks: {report['summary']['passed_checks']}")
        print(f"⚠️  Total issues: {report['summary']['total_issues']}")
        print(f"🔴 High severity: {report['summary']['high_severity']}")
        print(f"🟡 Medium severity: {report['summary']['medium_severity']}")
        print(f"🔵 Low severity: {report['summary']['low_severity']}")
        print(f"💡 Recommendations: {report['summary']['recommendations']}")
        
        if self.issues:
            print(f"\n🚨 SECURITY ISSUES")
            print(f"==================")
            for issue in self.issues:
                severity_emoji = {"HIGH": "🔴", "MEDIUM": "🟡", "LOW": "🔵"}
                print(f"{severity_emoji[issue['severity']]} [{issue['severity']}] {issue['category']}: {issue['description']}")
                if issue['file']:
                    print(f"   📁 File: {issue['file']}")
                    
        if self.passed_checks:
            print(f"\n✅ PASSED CHECKS")
            print(f"================")
            for check in self.passed_checks:
                print(f"✓ {check}")
                
        if self.recommendations:
            print(f"\n💡 RECOMMENDATIONS")
            print(f"==================")
            for rec in self.recommendations:
                print(f"💡 {rec['category']}: {rec['description']}")
                
        # Security score calculation
        total_checks = len(self.passed_checks) + len(self.issues)
        if total_checks > 0:
            score = (len(self.passed_checks) / total_checks) * 100
            print(f"\n🎯 SECURITY SCORE: {score:.1f}/100")
            
            if score >= 90:
                print("🟢 Excellent security posture!")
            elif score >= 75:
                print("🟡 Good security, minor issues to address")
            elif score >= 50:
                print("🟠 Moderate security, several issues need attention")
            else:
                print("🔴 Poor security, immediate action required")
        
        return report

def main():
    """Main function"""
    if len(sys.argv) > 1:
        project_root = sys.argv[1]
    else:
        project_root = os.path.dirname(os.path.abspath(__file__))
        
    auditor = SecurityAuditor(project_root)
    report = auditor.run_audit()
    
    # Save report to file
    report_path = Path(project_root) / 'security_audit_report.json'
    with open(report_path, 'w') as f:
        json.dump(report, f, indent=2)
        
    print(f"\n📄 Full report saved to: {report_path}")
    
    # Exit with error code if high severity issues found
    high_severity_count = report['summary']['high_severity']
    if high_severity_count > 0:
        print(f"\n❌ Audit failed: {high_severity_count} high severity issue(s) found")
        sys.exit(1)
    else:
        print(f"\n✅ Security audit passed!")
        sys.exit(0)

if __name__ == '__main__':
    main()
