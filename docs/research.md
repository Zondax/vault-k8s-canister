---
title: "Research"
sidebar_position: 4
---

# ICP Vault 🔑 :lock: Initial Research 

## :space_invader: Problem Definition
Customers face challenges in the secret management market due to complex deployment processes, high costs, vendor lock-in, security concerns in centralized systems, and limitations in customization and integration options. A new solution is needed to provide a user-friendly, secure, and cost-effective secret management experience that addresses these pain points while offering flexibility and adaptability to varying user needs.

From the point of view of customers, the current secret management market has several pain points and challenges that need to be addressed:


* **Security concerns**: In a centralized secret management system, customers may have concerns about data privacy and security, as there is a central point of failure and potential for data breaches.
* **Complex deployment and setup**: Customers often struggle with the time-consuming and complicated process of setting up and configuring secret management solutions, especially for organizations and enterprises.
* **High costs and vendor lock-in**: Some secret management solutions come with high costs, including licensing, maintenance, and professional services fees. Additionally, customers may face challenges in migrating to other solutions due to vendor lock-in.
* **Limited integrations:** Many secret management solutions have limited compatibility with other platforms, tools, and services, hindering seamless workflows and creating inefficiencies for users.
* **Inconsistencies in user experience**: Users may face inconsistencies in user experience across different platforms (desktop, web, mobile), leading to confusion and inefficiencies in secret management tasks.
* **Insufficient support and community resources:** Users often encounter inadequate customer support, especially for free users or specific user groups, and may find limited community resources for troubleshooting and learning.
* **Customization limitations:** Organizations with unique requirements may face difficulties customizing existing secret management solutions to meet their specific needs, leading to suboptimal implementations.

Zondax aims to develop a decentralized secret management product utilizing Internet Computer (ICP) technology to compete with existing solutions such as 1Password, Doppler, and Hashicorp Vault (among others). 

The product will offer a secure, decentralized infrastructure that eliminates the need for locally stored secrets, with the potential to evolve into a tokenized DAO. it will prioritize user experience, addressing common complaints about the complexity of AWS/GCP and the poor user experience of HCP.

## :chart_with_upwards_trend: Market Trends

* **Cloud adoption:** As more organizations migrate to the cloud, there is an increasing demand for cloud-native secret management solutions that integrate seamlessly with cloud platforms, such as AWS, GCP, and Azure.
* **Automation and DevOps integration:** With the rise of DevOps and CI/CD practices, there is a growing need for secret management solutions that can integrate with various DevOps tools and automate secret provisioning, rotation, and access management.
* **Zero-trust security:** Organizations are adopting zero-trust security models, which assume no inherent trust in any user or system. This trend drives the demand for secret management solutions that provide granular access controls, multi-factor authentication, and encryption.
* **Compliance and regulatory requirements:** As data privacy regulations like GDPR, CCPA, and HIPAA become more stringent, organizations need secret management solutions that help maintain compliance by providing access controls, audit logs, and data protection features.
* **Containerization and microservices:** The increasing adoption of containerization and microservices architectures leads to the need for secret management solutions that can securely manage secrets across distributed environments and dynamically provision secrets for containerized applications.
* **Open-source and community-driven solutions:** There is a growing interest in open-source secret management solutions that offer transparency, flexibility, and community-driven development, allowing organizations to customize and adapt the solution to their specific needs.
* **Multi-cloud and hybrid environments:** With organizations increasingly adopting multi-cloud and hybrid environments, there is a need for secret management solutions that can manage secrets across various platforms and on-premises systems seamlessly.
* **Decentralization and blockchain technology:** Emerging trends in decentralization and blockchain technology offer new opportunities for developing secure, transparent, and resilient secret management platforms that leverage the benefits of distributed systems.
* **Passwordless authentication and Passkey:** As organizations move towards more secure and user-friendly authentication methods, passwordless authentication and Passkey technology are gaining popularity. These methods, such as biometric authentication or single-use authentication codes, reduce the reliance on traditional passwords and the risks associated with password reuse or weak passwords.
* **Machine learning and AI-based security:** The integration of machine learning and artificial intelligence in secret management solutions can help detect anomalies, identify potential threats, and automate security responses, enhancing the overall security posture of the system.
* **Just-In-Time (JIT) access:** Organizations are adopting JIT access policies, granting users temporary access to secrets only when needed and automatically revoking access once the task is completed. This minimizes the risk of unauthorized access and reduces the attack surface.
* **Secrets-as-a-Service “SaaS”? :** The increasing demand for on-demand, scalable, and managed secret management solutions has led to the rise of Secrets-as-a-Service offerings. These services provide a cloud-based, fully managed secret management platform, reducing the need for organizations to maintain their own secret management infrastructure.
* **Unified secret management platforms:** As organizations manage various types of secrets, such as API keys, passwords, certificates, and encryption keys, there is a growing trend towards unified secret management platforms that can manage all secret types within a single, centralized solution.
* **Increased focus on user experience:** Secret management solutions are prioritizing user experience by designing intuitive interfaces, simplifying workflows, and offering seamless integrations with commonly used tools and platforms. This makes it easier for both technical and non-technical users to manage secrets effectively.

## :japanese_ogre: Competitor Analysis

We have identified and analyzed  the following competitors

1. Thycotic Secret Server
1. CyberArk Conjur 
1. Password Vault
1. ManageEngine Password Manager Pro
1. Keeper Security
1. LastPass Enterprise
1. Bitwarden Enterprise
1. Passbolt
1. TeamPassword
1. HashiCorp Vault
1. 1Password
1. Doppler
1. Akeyless


### Thycotic Secret Server:
   - Pros:
     - Comprehensive privileged access management solution.
     - On-premises and cloud deployment options.
     - Role-based access controls with granular permissions.
     - Integrations with various platforms, including Active Directory and SIEM tools.
     - Supports automated password rotation and customizable password policies.
   - Cons:
     - Steeper learning curve compared to some competitors.
     - Pricing can be high, especially for smaller businesses.
     - Lack of a mobile app for managing secrets on-the-go.


### CyberArk Enterprise Password Vault:
   - Pros:
     - Industry-leading solution with a strong focus on security and compliance.
     - Centralized management of privileged accounts, SSH keys, and certificates.
     - Supports automated password rotation and secure storage of old secrets.
     - Integrations with various platforms, including SIEM tools, and offers a robust API.
     - Provides detailed audit trails and reporting capabilities.
   - Cons:
     - Complex deployment and setup process.
     - High cost of ownership, including licensing, maintenance, and professional services fees.
     - The user interface can be less intuitive compared to some competitors.

### ManageEngine Password Manager Pro:
   - Pros:
     - Offers a wide range of features, including password management, access control, and auditing.
     - Integrations with various platforms, including Active Directory and LDAP.
     - Supports two-factor authentication and role-based access controls.
     - Offers both on-premises and cloud deployment options.
     - More affordable pricing compared to some competitors.
   - Cons:
     - The user interface could be more intuitive and visually appealing.
     - Limited mobile app functionality.
     - Customer support could be more responsive.

### Keeper Security:
   - Pros:
     - User-friendly interface and straightforward setup process.
     - Strong focus on security with zero-knowledge architecture and AES-256 encryption.
     - Supports multi-factor authentication, including biometrics and hardware tokens.
     - Offers a wide range of integrations, including SSO and popular web browsers.
     - Mobile applications available for iOS and Android devices.
   - Cons:
     - Limited support for automated secret rotation.
     - The granular access control options could be more comprehensive.
     - Pricing can be high for some businesses, especially those with a large number of

### LastPass Enterprise:
   - Pros:
     - User-friendly interface and easy setup.
     - Integrations with various platforms, including SSO, Active Directory, and web browsers.
     - Supports multi-factor authentication, including biometrics and hardware tokens.
     - Offers a shared folder feature for team collaboration.
     - Mobile applications available for iOS and Android devices.
   - Cons:
     - Limited support for automated secret rotation.
     - Granular access control options could be more comprehensive.
     - Some users have reported occasional sync issues and slow customer support response times.

### Bitwarden Enterprise:
   - Pros:
     - Open-source solution with a strong focus on security and transparency.
     - Competitive pricing, especially for smaller businesses.
     - Offers a wide range of integrations, including SSO and popular web browsers.
     - Supports multi-factor authentication with various methods.
     - Mobile applications available for iOS and Android devices.
   - Cons:
     - Limited support for automated secret rotation.
     - User interface may not be as polished as some competitors.
     - Smaller community and support resources compared to more established solutions.

### Passbolt:
   - Pros:
     - Open-source solution with a focus on transparency and collaboration.
     - Can be self-hosted or used as a cloud service.
     - Offers browser extensions and a command-line interface for easy access.
     - Supports GPG key-based authentication for increased security.
     - Competitive pricing, especially for small to medium-sized businesses.
   - Cons:
     - Limited integrations with other platforms and tools.
     - No mobile app available.
     - Lacks some advanced features, such as automated secret rotation and detailed access control options.

### TeamPassword:
   - Pros:
     - Simple and straightforward interface for managing shared passwords.
     - Integrations with popular web browsers for easy access.
     - Supports two-factor authentication for added security.
     - Offers a shared folder feature for team collaboration.
     - Suitable for small to medium-sized businesses.
   - Cons:
     - Limited support for automated secret rotation.
     - Lacks advanced features and integrations found in more comprehensive solutions.
     - No mobile app available.

### HashiCorp Vault:
   - Pros:
     - Open-source solution with a strong focus on security and scalability.
     - Unified secret management for various types of secrets, such as API keys, passwords, and certificates.
     - Supports dynamic secrets and automated secret rotation.
     - Offers a wide range of integrations, including cloud providers and popular DevOps tools.
     - Provides detailed audit logs and reporting capabilities.
   - Cons:
     - Steep learning curve and complex setup process.
     - The user interface could be more intuitive and user-friendly.
     - Limited support for some features, such as granular access controls, in the open-source version.

### 1Password:
   - Pros:
     - User-friendly interface and easy setup.
     - Strong focus on security with zero-knowledge architecture and AES-256 encryption.
     - Supports multi-factor authentication, including biometrics and hardware tokens.
     - Offers a wide range of integrations, including SSO and popular web browsers.
     - Mobile applications available for iOS and Android devices.
   - Cons:
     - Limited support for automated secret rotation.
     - Granular access control options could be more comprehensive.
     - Pricing can be high for some businesses, especially those with a large number of users.

### Doppler:
   - Pros:
     - Simple and straightforward interface for managing secrets and environment variables.
     - Integrations with various platforms, including CI/CD tools, cloud providers, and container orchestrators.
     - Offers real-time secret syncing and versioning.
     - Supports access controls and audit logs for compliance.
     - Competitive pricing, especially for small to medium-sized businesses.
   - Cons:
     - Limited support for some advanced features, such as automated secret rotation.
     - No mobile app available.
     - Lacks some integrations found in more comprehensive solutions.

### Akeyless:
   - Pros:
     - Offers centralized management of secrets, keys, and tokens.
     - Supports multi-cloud and hybrid environments.
     - Provides advanced access control options, including role-based access control and just-in-time access.
     - Offers integrations with popular DevOps tools and platforms.
     - Provides detailed audit logs and reporting capabilities.
   - Cons:
     - The user interface could be more intuitive and visually appealing.
     - Limited mobile app functionality.
     - Pricing can be high for smaller businesses.
