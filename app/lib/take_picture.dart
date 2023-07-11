import 'package:camera/camera.dart';
import 'package:flutter/material.dart';

import 'api_service.dart';

// A screen that allows users to take a picture using a given camera.
class TakePictureScreen extends StatefulWidget {
  const TakePictureScreen({super.key, required this.cameras});
  final List<CameraDescription>? cameras;

  @override
  TakePictureScreenState createState() => TakePictureScreenState();
}

class TakePictureScreenState extends State<TakePictureScreen> {
  late CameraController _controller;
  Future<void>? _initializeControllerFuture;
  late bool uploading = false;
  double _minAvailableZoom = 1.0;
  double _maxAvailableZoom = 1.0;
  double _currentZoomLevel = 1.0;
  bool _isRearCameraSelected = true;
  bool _isImageUploading = false;

  @override
  void initState() {
    super.initState();
    initCamera(widget.cameras![0]);
  }

  @override
  void dispose() {
    // Dispose of the controller when the widget is disposed.
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.black,
      appBar: AppBar(title: const Text('Take a picture')),
      body: _controller.value.isInitialized
          ? Column(
              children: [
                AspectRatio(
                    aspectRatio: 1 / _controller.value.aspectRatio,
                    child: Stack(children: [
                      CameraPreview(_controller),
                      Padding(
                          padding: const EdgeInsets.fromLTRB(
                            16.0,
                            8.0,
                            16.0,
                            8.0,
                          ),
                          child: Column(
                              mainAxisAlignment: MainAxisAlignment.end,
                              children: [
                                Row(
                                  children: [
                                    Expanded(
                                      child: Slider(
                                        value: _currentZoomLevel,
                                        min: _minAvailableZoom,
                                        max: _maxAvailableZoom,
                                        activeColor: Colors.white,
                                        inactiveColor: Colors.white30,
                                        onChanged: (value) async {
                                          setState(() {
                                            _currentZoomLevel = value;
                                          });
                                          await _controller.setZoomLevel(value);
                                        },
                                      ),
                                    ),
                                    Padding(
                                      padding:
                                          const EdgeInsets.only(right: 8.0),
                                      child: Container(
                                        decoration: BoxDecoration(
                                          color: Colors.black87,
                                          borderRadius:
                                              BorderRadius.circular(10.0),
                                        ),
                                        child: Padding(
                                          padding: const EdgeInsets.all(8.0),
                                          child: Text(
                                            _currentZoomLevel
                                                    .toStringAsFixed(1) +
                                                'x',
                                            style:
                                                TextStyle(color: Colors.white),
                                          ),
                                        ),
                                      ),
                                    ),
                                  ],
                                ),
                                Row(
                                  mainAxisAlignment: MainAxisAlignment.center,
                                  children: [
                                    //------ changement d'objectif
                                    widget.cameras!.length < 2
                                        ? Container()
                                        : InkWell(
                                            onTap: () {
                                              onNewCameraSelected(
                                                  widget.cameras![
                                                      _isRearCameraSelected
                                                          ? 1
                                                          : 0]);
                                              setState(() {
                                                _isRearCameraSelected =
                                                    !_isRearCameraSelected;
                                              });
                                            },
                                            child: Stack(
                                              alignment: Alignment.center,
                                              children: [
                                                Icon(
                                                  Icons.circle,
                                                  color: Colors.black38,
                                                  size: 60,
                                                ),
                                                Icon(
                                                  _isRearCameraSelected
                                                      ? Icons.camera_front
                                                      : Icons.camera_rear,
                                                  color: Colors.white,
                                                  size: 30,
                                                ),
                                              ],
                                            ),
                                          ),
                                    //------ déclencheur
                                    _isImageUploading
                                        ? CircularProgressIndicator(
                                            value: null,
                                            semanticsLabel: 'Image uploading',
                                          )
                                        : InkWell(
                                            onTap: takePicture,
                                            child: Stack(
                                              alignment: Alignment.center,
                                              children: [
                                                Icon(
                                                  Icons.circle,
                                                  color: Colors.white38,
                                                  size: 80,
                                                ),
                                                Icon(
                                                  Icons.circle,
                                                  color: Colors.white,
                                                  size: 65,
                                                ),
                                              ],
                                            ),
                                          )
                                  ],
                                )
                              ]))
                    ]))
              ],
            )
          : const Center(child: CircularProgressIndicator()),
/*       floatingActionButton: uploading
          ? const Center(child: CircularProgressIndicator())
          : FloatingActionButton(
              // Provide an onPressed callback.
              onPressed: () async {
                // Take the Picture in a try / catch block. If anything goes wrong,
                // catch the error.
                try {
                  // Ensure that the camera is initialized.
                  await _initializeControllerFuture;

                  // Attempt to take a picture and get the file `image`
                  // where it was saved.
                  final image = await _controller.takePicture();

                  setState(() {
                    uploading = true;
                  });

                  final imageUrl = await uploadImage(image);

                  if (!mounted) return;
                  Navigator.pop(context, imageUrl);
                } catch (e) {
                  // If an error occurs, log the error to the console.
                  print(e);
                }
              },
              child: const Icon(Icons.camera_alt),
            ),
 */
    );
  }

  Future initCamera(CameraDescription cameraDescription) async {
    // create a CameraController
    _controller = CameraController(cameraDescription, ResolutionPreset.medium);
    // Next, initialize the controller. This returns a Future.
    try {
      await _controller.initialize().then((_) {
        if (!mounted) return;
        _controller
            .getMaxZoomLevel()
            .then((value) => _maxAvailableZoom = value);

        _controller
            .getMinZoomLevel()
            .then((value) => _minAvailableZoom = value);

        setState(() {});
      });
    } on CameraException catch (e) {
      debugPrint("camera error $e");
    }
  }

  takePicture() async {
    // Take the Picture in a try / catch block. If anything goes wrong,
    // catch the error.

    try {
      // Ensure that the camera is initialized.
      await _initializeControllerFuture;

      // Attempt to take a picture and get the file `image`
      // where it was saved.
      final image = await _controller.takePicture();

      setState(() {
        uploading = true;
      });

      final imageUrl = await uploadImage(image);

      if (!mounted) return;
      Navigator.pop(context, imageUrl);
    } catch (e) {
      // If an error occurs, log the error to the console.
      print(e);
    }
  }

  void onNewCameraSelected(CameraDescription cameraDescription) async {
    final previousCameraController = _controller;

    final CameraController cameraController = CameraController(
      cameraDescription,
      ResolutionPreset.medium,
      imageFormatGroup: ImageFormatGroup.jpeg,
    );

    try {
//      await previousCameraController.dispose();
    } catch (e) {}

    if (mounted) {
      setState(() {
        _controller = cameraController;
      });
    }

    // Update UI if controller updated
    cameraController.addListener(() {
      if (mounted) setState(() {});
    });

    try {
      await cameraController.initialize();
      await Future.wait([
        cameraController
            .getMaxZoomLevel()
            .then((value) => _maxAvailableZoom = value),
        cameraController
            .getMinZoomLevel()
            .then((value) => _minAvailableZoom = value),
      ]);
    } on CameraException catch (e) {
      print('Error initializing camera: $e');
    }

    if (mounted) {
      setState(() {});
    }
  }

  Future<String> uploadImage(XFile image) async {
    setState(() {
      _isImageUploading = true;
    });
    final imageBytes = await image.readAsBytes();
    String result = await ApiService().postPicture(imageBytes);
    _isImageUploading = false;
    return result;
  }
}
