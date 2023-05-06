---


> This repository is a mirror of GitLab repository.


----

**EZDEB** is a CLI application that helps users manage debian packages from 3rd party websites and GitHub repositories. 

All software packages are not available on distribution's official repositories which is why ezdeb was created, it aims to maintain an up-to-date list of
applications (.deb packages) distributed through GitHub and websites.

## Features:    
- Package listing, information and searching  
  - List packages  
    - View installed packages only  
    - View held packages only  
  - View package information  
  - Search packages
  - View application version
- Package installation and uninstallation
  - Install package(s)
  - Uninstall package(s)
- Package updating, syncing and config file management
  - Update packages
    - Only check for updates
  - Sync package list
- Action logs and temporary files management
  - View logs
    - View logs for specific action
  - Clear logs
  - Clean temporary files
- Hold, unhold packages
  - Hold packages
  - Unhold packages
- CLI application information
  - View CLI usage help and individual commands help
  - View CLI application version
  
## Screenshots

## Building ezdeb from source:

Ensure you have a properly set-up GO  installation on your system.

Open a terminal window

Execute the following command:
```
git clone https://github.com/Charlie-117/ezdeb
```

Change into the directory
```
cd ezdeb
```

Build the binary for the application
```
GOOS=linux GOARCH=amd64 go build -o ezdeb-linux-amd64 main.go
```

That's it! You have successfully built the binary application for ezdeb from source.
