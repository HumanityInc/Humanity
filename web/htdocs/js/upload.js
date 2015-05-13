function DestroyUploadImg() {
	if (window._cur_uploader) {
		window._cur_uploader.destroy();
		window._cur_uploader = null;
	}
}

function InitUploadFile(upload_url, container, browse_button, drop_element, cb_img_url, cb_error, cb_start, cb_progress)
{
	DestroyUploadImg();

	var uploader = new plupload.Uploader({
		runtimes: 'html5,flash,silverlight',
		browse_button: browse_button,
		container: container,
		drop_element: drop_element,
		multi_selection: false,
		multipart_params: {format: "plain"},
		urlstream_upload: true,
		file_data_name: 'file',
		max_file_size: '16mb',
		url: upload_url,
		flash_swf_url: '/js/plupload/Moxie.swf',
		silverlight_xap_url: '/js/plupload/Moxie.xap',
		filters: [{
			title: "Images",
			extensions: "png,jpg,jpeg,bmp,tif,tiff"
		}]
	});

	window._cur_uploader = uploader;

	uploader.bind('Init', function(up, params) {
		console.log(params.runtime);
	});

	uploader.init();

	uploader.bind('FilesAdded', function(up, files) {
		up.refresh();
		uploader.start();
		if (cb_start) cb_start(files);
	});

	uploader.bind('UploadProgress', function(up, file) {
		if (cb_progress) cb_progress(file);
	});

	uploader.bind('Error', function(up, err) {
		// console.log("Error: ", err.message);
		up.refresh();
	});

	uploader.bind('FileUploaded', function(up, file, res) {

		if (res.status == 200) {

			try {
				var response = $.parseJSON(res.response);
				if (cb_img_url) cb_img_url(response.filelink, response.cover, file);
			}
			catch(e) {}
		}
		else {
			if (cb_error) cb_error(file);
			alert( "Произошла ошибка при загрузки картинки. Пожалуйста, попробуйте позже." );
			// console.log(res, file);
		}
	});
}
