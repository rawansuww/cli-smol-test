simple CLI to load test ENDUSER ditto API with pre-prod tokens

plz add config.yaml and add tokens in this format:
tokens:
- Bearer xxxxx

then


go build -o ./<name>

./<name> stresstest <URL> <INT>