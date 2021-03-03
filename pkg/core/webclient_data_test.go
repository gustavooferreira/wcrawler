package core_test

const htmlBody1 = `
<!DOCTYPE html>
<html>
<head>
<title>Title</title>
</head>
<body>
<a href="http://www.example.com/file.html#frag1">link1</p>
<a href="/path/to/file999">link1</p>
<a href="path/to/file2#frag123">link1</p>
</body>
</html>
`

const htmlBody2 = `
<!DOCTYPE html>
<html>
<head>
<title>Title</title>
<base href="http://www.example.com/base/path/to/dir/" target="_blank" />
</head>
<body>
<a href="/path/to/file1">link1</p>
<a href="relative/file2">link1</p>
</body>
</html>
`
