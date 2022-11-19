# heatcontroller

Create executable : go build gpio.go heatcontroller.go
Run it as a system service : sudo systemctl start gpio.service
System service status : sudo systemctl status gpio.service
Stop system service : sudo systemctl stop gpio.service



Remarks on Git:

Using git with write access:

1.) Make sure git is installed
2.) Goto github to your repo, copy the "clone" link
3.) Go to your local folder where this clone should be pulled to
4.) git clone <COPIED_CLONE_LINK>  (e.g. git clone https://github.com/wir33658/heatcontroller.git)
5.) git pull
6.) Do changes in your code ...
7.) git commit -a -m "minor"  (-a  add  -m  comment)
8.) Set User for this repo : git config user.email "robert.weissmann@web.de"  and  git config user.name "wir33658
9.) Till here is only reading and writing to the local repo. Now we need a token to get write access.
    Goto : https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token
    and follow the instructions -> choose "All repositories" and "Contents : Read and Write" -> create token
10.) Copy the created token and : git remote set-url origin https://wir33658:<COPIED_TOKEN>@github.com/wir33658/heatcontroller.git
11.) Now you have write access : git push  is now possible.


Remarks on server.go:

To work with TLS (https) a certificate needs to be created. Check also https://www.golinuxcloud.com/golang-http/

Run:

- openssl genrsa -out ca.key 4096
- openssl req -new -x509 -days 365 -key ca.key -out cacert.pem -subj "/C=IN/ST=NSW/L=Oslo/O=GoLinuxCloud/OU=Org/CN=RootCA"

-> Output : ca.key   and    cacert.pem

- Create file : server_cert_ext.cnf (see https://www.golinuxcloud.com/golang-http/)
- openssl genrsa -out server.key 4096
- openssl req -new -key server.key -out server.csr -subj "/C=IN/ST=NSW/L=Oslo/O=GoLinuxCloud/OU=Org/CN=imac"
- openssl x509 -req -in server.csr  -CA cacert.pem -CAkey ca.key -out server.crt -CAcreateserial -days 365 -sha256 -extfile server_cert_ext.cnf

-> Output : server.crt  and   server.csr  and   server.key

- cacert.pem and server.crt need to be bundled up since the parameter for the http.ListenAndServeTLS(...) function :
  cp server.crt certbundle.pem
  cat cacert.pem >> certbundle.pem
  Check it : cat certbundle.pem
