# fabric ca

## 参考资料

- https://medium.com/@_mtnieto/how-to-set-up-a-hyperledger-fabric-infrastructure-using-certificate-authorities-a00a6f01b35e
- https://stackoverflow.com/questions/58543416/run-intermediate-ca-with-tls-enabled-connection-to-root-ca-refused


Error: POST failure of request: POST http://ca.root:7054/enroll
{"hosts":["admin"],
"certificate_request":"-----BEGIN CERTIFICATE REQUEST-----\nMIIBPDCB5AIBADBfMQswCQYDVQQGEwJVUzEXMBUGA1UECBMOTm9ydGggQ2Fyb2xp\nbmExFDASBgNVBAoTC0h5cGVybGVkZ2VyMQ8wDQYDVQQLEwZGYWJyaWMxEDAOBgNV\nBAMTB2FkbWluQ0EwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAASJRk72ffXSDXGQ\noA8RVkST/jJaWRiJx5RvvMsAP7ir4CIvWtQ7iJeOurW0LU69sqWlaK34pAOf9nGk\n9b07bGW0oCMwIQYJKoZIhvcNAQkOMRQwEjAQBgNVHREECTAHggVhZG1pbjAKBggq\nhkjOPQQDAgNHADBEAiBKFOTX9/DOfZhgp8Vnkx0DDMTmk7QcbOZSgV7Ab7Iu5wIg\nGZPn4NL62ga45HuKAXXKrPRhEIC98FcHSuolBmqwxjc=\n
-----END CERTIFICATE REQUEST-----\n","profile":"","crl_override":"","label":"","NotBefore":"0001-01-01T00:00:00Z","NotAfter":"0001-01-01T00:00:00Z","CAName":""}: 
Post http://ca.root:7054/enroll: dial tcp: lookup ca.root on 180.76.76.76:53: no such host