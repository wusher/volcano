

this is a project that is a go commandline. 

the two main ways it should be run are 


where applicable or when guidance is lacking  it should follow the patterns set forward by https://github.com/wusher/medusa-ssg 



## generate static site 

`volcano FOLDER_NAME -o OUTPUT_FOLER`

you can improve the flags 


running it this will take a folder of markdown files and generate a static website  



## serve a static site 

`volcano -s -p 3000 FOLDER_NAME `

-p == port and is optional 


## layouts / css 


all layouts and css should be including in volcano. it hsould not need any depenancies 


for now everything should be simple. 

### title 


the flag --title="TEXT" will pass the site's title. default value is My Site 


### layout 

use the markdown files and the folders ( ignoring empty ones ) to build a tree layout on the left hand side 

in mobile view, the nav is a drawer slideout 


### content the right hand side of the page is 

the markdown renders as html and styled using tailwind typography 



###  colors 

it should be shades of black and white  
there should be a toggle for light and dark mode and default to browers preference 




## dev notes 

use the folder example as test imput 
code should have 95% test coverage 
code should have e2e tests 
code should be linted 
