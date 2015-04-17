# fdu.go

This is my first ever go program.

As the ruby and python version, this is a usefull tool:

 $ fdu dir1 dir2/* dir3 dir4/* file4

will give you an estimate of the size of the directories/files, like du -csh.

However, each directory will have a .SIZE file that cache the result.
Launching this regularly on backup or log servers for instance will help 
knowing which directories are the bigest.

So that the second time one will launch the command, the answer will be instant.

Lots could be better.

Fill free to use, modify or rewrite from scratch ;-)


