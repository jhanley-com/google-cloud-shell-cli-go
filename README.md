# google-cloud-shell-cli-go
Repository for my article on Google Cloud Shell CLI in Go

<pre>
Usage: cloudinfo [command]
  cloudinfo                            - display Cloud Shell information
  cloudinfo info                       - display Cloud Shell information
  cloudinfo putty                      - connect to Cloud Shell with Putty
  cloudinfo exec "command"             - Execute remote command on Cloud Shell
  cloudinfo upload src_file dst_file   - Upload local file to Cloud Shell
  cloudinfo download src_file dst_file - Download from Cloud Shell to local file

--debug - Turn on debug output
--auth  - (re)Authenticate ignoring user_credentials.json
--login - Specify an email address as a login hint
</pre>
