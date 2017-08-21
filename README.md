# contactkey [![Build Status](https://travis-ci.org/remyLemeunier/contactkey.svg?branch=master)](https://travis-ci.org/remyLemeunier/contactkey)

Contactkey is a Go command and its aim is to deploy services on different environments.

## Deployment Flow
![this](https://docs.google.com/drawings/d/1N7mgky_Dq3KWrT_gRxR4iwxGjCDY6rbgc455mJgEMtA/pub?w=594&h=1155)
# Deployers, VCS, Repository Manager, Hooks, Lock System
### Deployers
 * GGN 
 * kubernetes (to do) 
### Version Control System
 * Stash 
 * Github (to do)
### Repository Manager 
 * Nexus
### Hooks
 * Slack
 * ExecCommand
 * NewRelic (to do)
 * MailSender (to do)
### Lock System 
 * FileLock
 * etcd (to do)
 * Redis (to do)
 
## Commands
```bash
cck deploy environment service    Deploy the service in an environment
cck diff environment service      Diff between what's currently deployed and what's going to be deployed (VCS)
cck list environment service      List versions of the service in an environment
cck rollback environment service  Rollback the service in an environment
```
## Global configuration file for the cck command.
The configuration file is located at ~/.contackey/config.yml
```yaml
workPath: /tmp/manifests   # Location of services manifest

screenMandatory: true      # Check if the user is launching cck in a screen/tmux (not mandatory)

globalEnvironments:        # Define the cck environment. It can be anything.
  - preprod                # It will be used as the cck environment for 
  - prod                   # the command line.

deployers:                 # Define the various deployers used in service manifest.
  ggn:                     # Currently we have only ggn supported .
    vcsRegexp: -v(.+)      # Extract the vcs sha1 from pod version. (Not mandatory)
    workPath: /tmp       
    environments:          # Link between ggn environment and cck environment created above. 
      preprod: staging    
      prod:    production 

versionControlSystem:      # Define various version control system used in service manifest.
  stash:                   # Currently we have only Stash supported. 
    user:        user     
    password:    password 
    url:         url       
    sha1MaxSize: 2         # Cut a sha1 if it's too long. (E.g: abcd => ab) (Not mendatory) 

repositoryManager:         # Define various repository manager used in service manifest. 
  nexus:                   # Currently we have only Nexus supported. 
    url:        127.0.0.1  
    repository: repository 
    group:      group      

hooks:                     # Define various hooks used in service manifest. 
  slack:                   # Currently we have only Slack supported.
    url:   127.0.0.1      
    token: token          

lockSystem:                # Define a lock system in order to avoid multiple command launch. (Not mandatory)
  fileLock:                # Currently we only have the lock by file.
    filePath: /tmp         # Path where the is going to be written.

```
## Configuration by service (manifest)
One file per service. They must be located in the workPath defined above.
The file name will be service name in the cck command.
```yaml
deployment:                # Define the deployment type we are going to use
  ggn:                     # for this service.
    pod: pod-webhooks      # podName used in ggn
    service: webhooks      # service used in ggn
versionControlSystem:      
  branch: master           # The default "stable" branch we usually deploy
  stash:                   # Version control system used for the service. (Only one) 
    repository: repository 
    project: project       
repositoryManager:         
  nexus:                   # Repository manager used for the service. (Only one)
    artifact: pod-webhooks 
hooks:                     # Hooks we are going to call before and after.
  slack:                   # the deployment. (You can have several one)
    channel:   channel
    StopOnError: false     # If an error occurs stop the deployment process (not mandatory default false)
  execCommand:             # Execute a command before and after the deployment process.
    list:
      - "cd /tmp"
      - "ls"
    StopOnError: true      # If an error occurs stop the deployment process (not mandatory default false)
```