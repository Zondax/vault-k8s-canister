PodA
    Read SecretName1234

PodB
    Apply SecretName1234

Secret  Name1234                CRD
    Type: Postgres
    Rotate, etc.....

----------------------------------------------------  ICP
Secret  Name1234
    Type: Postgres
    Rotate, etc.....

    Will create and rotate:
    * Postgres User
    * Postgres Password
