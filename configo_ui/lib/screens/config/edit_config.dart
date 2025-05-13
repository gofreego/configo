import 'package:flutter/material.dart';
import 'package:configo_ui/models/config/object.dart';
import 'package:configo_ui/widgets/form.dart';
import 'package:configo_ui/services/config/config.dart';
import 'package:configo_ui/widgets/error.dart';

class ConfigForm extends StatefulWidget {
  final String configKey;
  final Function() onCancel;
  const ConfigForm({super.key, required this.configKey, required this.onCancel});

  @override
  State<ConfigForm> createState() => _ConfigFormState();
}

class _ConfigFormState extends State<ConfigForm> {
  bool isLoading = true;
  String? error;
  List<ConfigObject> configs = [];
  void fetchConfig() async {
    setState(() {
      isLoading = true;
    });
    var res = await ConfigService().getConfig(widget.configKey);
    setState(() {
      isLoading = false;
      configs = res.data?.configs ?? [];
      error = res.error;
    });
  }

  void saveConfig() async {
    setState(() {
      isLoading = true;
    });
    var res = await ConfigService().updateConfig(widget.configKey, configs);
    setState(() {
      isLoading = false;
      error = res.error;
      if (res.error == null) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text("Configs updated successfully"),
            backgroundColor: Color.fromARGB(255, 42, 95, 44),
          ),
        );
      }
    });
  }

  @override
  void initState() {
    fetchConfig();
    super.initState();
  }

  //show confirm dialog on cancel and submit
  confirmDialog(String title, String content, Function() onConfirm) {
    showDialog(
      context: context,
      builder: (BuildContext context) {
        return AlertDialog(
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(15),
          ),
          title: Text(
            title,
            style: TextStyle(
              fontWeight: FontWeight.bold,
              color: Theme.of(context).primaryColor,
            ),
          ),
          content: Text(content, style: TextStyle(fontSize: 16)),
          actions: <Widget>[
            TextButton(
              onPressed: () {
                onConfirm();
                Navigator.of(context).pop();
              },
              child: const Text(
                'Yes',
                style: TextStyle(
                  color: Color.fromARGB(255, 43, 100, 45),
                  fontWeight: FontWeight.bold,
                ),
              ),
            ),
            TextButton(
              onPressed: () {
                Navigator.of(context).pop();
              },
              child: const Text(
                'No',
                style: TextStyle(
                  color: Color.fromARGB(255, 146, 13, 4),
                  fontWeight: FontWeight.bold,
                ),
              ),
            ),
          ],
        );
      },
    );
  }

  onCancel() {
    widget.onCancel();
  }

  onSubmit() {
    confirmDialog(
      "Update Configs",
      "Are you sure you want to update the configs?",
      saveConfig,
    );
  }

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(16),
      child: Column(
        children: [
          isLoading
              ? const Center(
                child: Padding(
                  padding: EdgeInsets.all(8.0),
                  child: CircularProgressIndicator(),
                ),
              )
              : error != null
              ? CustomErrorWidget(errorMessage: error!, onRetry: fetchConfig)
              : Column(
                children: [
                  ConfigFormWidget(configs: configs),
                  const SizedBox(height: 10),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      ElevatedButton(
                        onPressed: onCancel,

                        child: const Text("Cancel"),
                      ),
                      const SizedBox(width: 10),
                      ElevatedButton(
                        onPressed: onSubmit,
                        style: ElevatedButton.styleFrom(
                          backgroundColor: Theme.of(context).primaryColor,
                        ),
                        child: const Text(
                          "Submit",
                          style: TextStyle(color: Colors.white),
                        ),
                      ),
                    ],
                  ),
                ],
              ),
        ],
      ),
    );
  }
}
