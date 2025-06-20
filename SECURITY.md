# Security Policy

## Supported Versions

Currently, we are actively maintaining and providing security updates for the following versions of Investify:

| Version | Supported          |
| ------- | ------------------ |
| 1.0.x   | :white_check_mark: |

## Reporting a Vulnerability

We take the security of Investify seriously. If you believe you've found a security vulnerability, please follow these steps:

1. **Do not disclose the vulnerability publicly**
2. **Email us** at [your-email@example.com] with details about the vulnerability
3. Include the following information:
   - Type of vulnerability
   - Steps to reproduce
   - Potential impact
   - Any suggested fixes (if known)

## Security Features

Investify implements the following security features:

- Input validation and sanitization for all user-provided data
- Content Security Policy (CSP) headers to prevent XSS attacks
- Secure storage of API keys through environment variables
- HTTPS support for secure communication
- Regular dependency updates to patch known vulnerabilities

## Best Practices for Users

When using Investify, please follow these security best practices:

1. Keep your API keys confidential and do not hardcode them in your files
2. Use environment variables for sensitive information
3. Regularly update your dependencies with `go get -u` and `pip install -r requirements.txt --upgrade`
4. Do not expose the application directly to the internet without proper authentication

## Acknowledgments

We would like to thank the following individuals who have responsibly disclosed security issues to us:

- Your name could be here!
