# appligen
Generates (*.pdf*-) job applications out of a latex template and a yaml-file.

# requirements
- ```latex```
- ```pdflatex```
- template file (```application.tex```) in the root directory
- data file (```data.yaml```) in the root directory
- ```figures``` folder in the root directory
  - ```signature.png```-file with your signature in the ```figures```-folder
 
# what it does
*appligen* creates subfolders for each item in the *applications*-list in the ```data.yaml```-file and generates a ```data.tex``` file with ```latex```-commands out of the dataset, that are used in the template and copies the folder ```figures``` and the ```application.tex``` (as ```template.tex```) into the subfolders. Afterwards it runs ```pdflatex template.tex``` in each subfolder to generate the final *.pdf*-file.

# why it was created
I was annoyed by the fact that every time I write an application I had to manually fiddle with my template. Now I just fill the *.yaml*-file with the necessary data, run the binary and I'm ready to go.

Now I can even create multiple applications with one run.
