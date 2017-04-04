# GoDeployServer

A little web server in Go to redeploy apps in GitHub push.

## Use

1. Modify the code to add your paths and your management scripts
2. Run `go build`
3. `sudo cp getLastCommit /usr/local/bin`
4. Run the service `systemctl start godeployserver`

## PostInfo

I'm using the PostInfo struct to process json payload of the GitHub webhook.
