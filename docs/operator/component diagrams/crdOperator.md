# CRD Operator

```mermaid
sequenceDiagram
    participant CRDOperator
    participant Kubernetes
    participant Canister
    participant CRDInstance

    CRDOperator->>Kubernetes: Listen for CRD Instance Creation/Update ✅

    alt CRD Created
        CRDOperator->>CRDInstance: Override Fields (e.g., secret = "") ✅
        note over CRDOperator, CRDInstance: Fields overridden

    else CRD Updated
        CRDOperator->>CRDInstance: Check podsRestartRequired Flag ✅
        note over CRDOperator, CRDInstance: podsRestartRequired = true

        alt Restart Required
            CRDOperator->>Canister: Get RO Consumers from State ✅
            note over CRDOperator, Canister: Retrieved RO Consumers

            alt Consumers Restarted Successfully
                CRDOperator->>Kubernetes: Restart RO Consumers ✅
                note over CRDOperator, Kubernetes: RO Consumers restarted successfully

                CRDOperator->>CRDInstance: Restore podsRestartRequired Flag to false ✅
                note over CRDOperator, CRDInstance: podsRestartRequired = false
            else Restart Failed
                CRDOperator->>Kubernetes: Log Restart Failure ✅
                note over CRDOperator, Kubernetes: RO Consumers restart failed
            end
        else No Restart Required
            note over CRDOperator, CRDInstance: podsRestartRequired = false
        end
    end
```

[✏️ Edit here](https://mermaid.live/edit#pako:eNqtlc1u2zAMx1-F0GkD3CxOnQ8baC_pMgz7CJBgl8EXzWJSobaUUXKxLOixb7Gn25NMtmPXXtKmBZqTI5J__ijS9I4lWiCLmMGfOaoEryRfE89iBe634WRlIjdcWZguruYbJG41HRo_5T-QFFo0RwK5ksYiHZX8qIzlLm2sKnMry9nl5YNsBJ8LEQUrTYUT1IEwJeRWavXu20Zwi_D3z30txtMyR-WCojo8zNLCiGB-i0RSIMwkpsLAG-ytex4YTAgtXEDMYva2SlKrKe3SahfX1vWgI7tX05W6QFUzYmqwhKzwnwc5vcbkBjZamAW6M7IL1z1JKGCW8vXL6Y4pXYClvOlLfZ17L6jdHqxHkPeNj-CDu7rFHKZamTxDMrAincHSNv1qizwG3Kgt0JLEW8fY1myTNs1vEu65XcwyTxI0ZpWn6bYb8fT8NZW36ziAf6KAjlhbhBo202E7RddpYYGnCR8fCqthxYthewnzM6akFO0KljNd39eMy_T_OTnxput1JzonfNWLdsiHSKhaB2UBX_WJaX-1O2tyuwfmMceacSncUt4VxzGz15hhzCL3KDjduBWk7pwfz61eblXCouJN9VhebpD9AmdRmcRjbtN-1zqrndxfFu3YLxYNwt4o9AfhcByOhoPJaDL02JZFZ8G4N_BH_Uk48QP_POgP7zz2u1To98Jg4AfOch6MJ8Ng5HsMhXRVf6k-IuW3pAZ5X1qqvHf_AMmOEn4)
