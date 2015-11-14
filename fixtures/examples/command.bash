// On UNIX
output = ls();
// On Windows
output = dir();
// Platform independent
output = readdir();

// Test existence
ex = exists("file.txt");
