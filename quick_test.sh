#!/bin/sh
rm -r envs
touch .env
cat <<EOF >.env
S3_BUCKET=YOURS3BUCKET
SECRET_KEY=YOURSECRETKEYGOESHERE
MALFORM=SECRETDONTTELL
EOF

go test -v envs_test.go envs.go crypto.go parser.go
