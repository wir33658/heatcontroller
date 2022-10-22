# heatcontroller



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

