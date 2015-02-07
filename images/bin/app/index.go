package app

func (c *Callback) wIndex() {

	form := ""

	if c.ip() == "127.0.0.1" {

		// test form

		form = `
			<html>
			<title>UPLOAD</title>
			<body>
				<form action="/upload" method="post" enctype="multipart/form-data">
					<label for="file">Filename:</label>
					<input type="file" name="file" id="file">
					<input type="submit" name="submit" value="Submit">
				</form>
			</body>
			</html>`
	}

	c.res.Write([]byte(form))
}
