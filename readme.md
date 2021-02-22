# Image Squeezer
Simple image optimization web app that takes user uploaded image files (PNG or JPEG) 
and compress it without any significant image quality loss. 
The width and height of the compressed image stays the same as initial image.
For some images the size might be reduced up to 75% of initial size.

Check it out here: http://imagesqueezer.com/

For image processing it uses MozJPEG and PNGQuant C libraries under the hood.
Requests handling and templates are implemented by standard Go packages.
[Gorilla sessions](https://github.com/gorilla/sessions) package to store cookies and sessions.

### Features
* Graceful server shutdown
* Processed image deletion using time intervals
* Continuous deployment using Jenkins and Docker Swarm
* Testing, linting and vetting during the deployment process (see Jenkinsfile)