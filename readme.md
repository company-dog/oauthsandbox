## Just a plain box here.

Centos7 | a few utils | PasswordAuthentication yes.

## setup keycloak

```
wget https://downloads.jboss.org/keycloak/8.0.1/keycloak-8.0.1.tar.gz

tar xvfz keycloak-8.0.1.tar.gz
cd keycloak-8.0.1/bin

sudo yum install java-1.8.0-openjdk-devel.x86_64


./standalone.sh -b 10.100.196.60
```

```
./add-user-keycloak.sh -u admin
```

