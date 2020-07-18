# reczip
Simple CLI tool to zip all files from a specified directory and its subdirectories. It puts an each file in a separate archive.

### Arguments:
+ --path  
    starting (root) directory
+ --ext  
    filter files by extension
+ --del=true  
    delete file after archiving

### Example:
reczip --path /home/user/text --ext .txt --del=true  
Archives all txt files in /home/user/text each to a single file archive {txt filename}.txt.zip, deletes source txt file after archiving.
