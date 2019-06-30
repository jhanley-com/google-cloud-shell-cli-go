# google-cloud-shell-cli-go
Repository for my article on Google Cloud Shell CLI in Go

https://www.jhanley.com

This program requires Google OAuth 2.0 Client Credentials. Go to the Google Cloud Console -> APIs & Services. Select Create credentials -> OAuth Client ID. Select "Other" for the Application Type.

Once you have OAuth 2.0 Client Credentials, edit the the file config.json to specify the full path to the credentials file.

The first time you execute this program, you will be prompted to authenticate with Google. These credentials are saved in the file user_credentials.json. In the module auth.go, I show how to store credentials and refresh the access token.

Note: This program is written for Windows. I have not yet ported it to Linux or Mac OS. Volunteers?

<pre>
Usage: cloudshell [command]
  cloudshell                            - display Cloud Shell information
  cloudshell info                       - display Cloud Shell information
  cloudshell putty                      - connect to Cloud Shell with Putty
  cloudshell exec "command"             - Execute remote command on Cloud Shell
  cloudshell upload src_file dst_file   - Upload local file to Cloud Shell
  cloudshell download src_file dst_file - Download from Cloud Shell to local file

--debug - Turn on debug output
--auth  - (re)Authenticate ignoring user_credentials.json
--login - Specify an email address as a login hint
</pre>

# Getting Started

This program supports Putty for the SSH client. Download https://www.putty.org/

You will need to setup the SSH keys for Cloud Shell. This requires the "alpha" version of the Google Cloud SDK:

1) Install the alpha components: <code>gcloud components install beta</code>
2) Update the Google Cloud SDK: <code>gcloud components update</code>
3) Create the SSH key pair: <code>gcloud alpha cloud-shell ssh --dry-run</code>

Install the Go dependencies:
<pre>
go get github.com/kirinlabs/HttpRequest
go get github.com/pkg/sftp
go get golang.org/x/crypto/ssh
</pre>

Build the program:
<pre>
go build -o cloudshell.exe
</pre>

## Examples
Run the program and display information about your Google Cloud Shell instance:
<pre>
cloudshell.exe info
</pre>

Launch Putty and connect to Cloud Shell:
<pre>
cloudshell.exe putty
</pre>

Upload a file to Cloud Shell:
The local file "local_file.txt" will be copied to the current working directory in Google Cloud Shell.
<pre>
cloudshell.exe upload local_file.txt remote_file.txt
</pre>

Quick file copy:
This command copies the local file "myfile.txt" to the Cloud Shell default working directory with the same file name.
<pre>
cloudshell.exe upload myfile.txt
</pre>

Copy a file to a specic location:
<pre>
cloudshell.exe upload local_file.txt /tmp/remote_file.txt
</pre>

What is the current Cloud Shell working directory?
<pre>
cloudshell.exe exec "pwd"
</pre>

Display the Cloud Shell current working directory files (directory listing):
<pre>
cloudshell.exe exec "ls -l"
</pre>

#### Note: The remote command must be enclosed in quotation marks
Remote commands that change the environment work but have no effect on the next command. You can combine commands in one session: <code>cloudshell exec "cd /home; cat testfile.txt"</code>
