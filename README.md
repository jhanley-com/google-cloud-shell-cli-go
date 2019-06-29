# google-cloud-shell-cli-go
Repository for my article on Google Cloud Shell CLI in Go

This program requires Google OAuth 2.0 Client Credentials. Go to the Google Cloud Console -> APIs & Services. Select Create credentials -> OAuth Client ID. Select "Other" for the Application Type.

Once you have OAuth 2.0 Client Credentials, edit the the file config.json to specify the full path to the credentials file.

The first time you execute this program, you will be prompted to authenticate with Google. These credentials are saved in the file user_credentials.json. In the module auth.go, I show how to store credentials and refresh the access token.

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
