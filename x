diff --git a/auth.go b/auth.go
index 73e0ee6..84ce9cc 100644
--- a/auth.go
+++ b/auth.go
@@ -569,7 +569,7 @@ func get_tokens() (string, string, error) {
 
 	auth_code := string(out)
 
-	return processAuthCode(secrets, auth_code)
+	return processAuthCode(secrets, auth_code, false)
 }
 
 func get_sa_tokens() (string, string, error) {
@@ -647,16 +647,20 @@ func manualAuthentication(secrets ClientSecrets, url string) (string, string, er
 
 	auth_code := strings.Replace(text, "\n", "", -1)
 
-	return processAuthCode(secrets, auth_code)
+	return processAuthCode(secrets, auth_code, true)
 }
 
-func processAuthCode(secrets ClientSecrets, auth_code string) (string, string, error) {
+func processAuthCode(secrets ClientSecrets, auth_code string, flag_oob bool) (string, string, error) {
 	//************************************************************
 	content := "client_id=" + secrets.Installed.ClientID
 	content += "&client_secret=" + secrets.Installed.ClientSecret
 	content += "&code=" + auth_code
 	content += "&grant_type=authorization_code"
-	content += "&redirect_uri=urn:ietf:wg:oauth:2.0:oob"
+	if flag_oob == false {
+		content += "&redirect_uri=http://localhost:9000"
+	} else {
+		content += "&redirect_uri=urn:ietf:wg:oauth:2.0:oob"
+	}
 	//************************************************************
 
 	endpoint := "https://www.googleapis.com/oauth2/v4/token"
