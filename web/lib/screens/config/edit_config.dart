import 'package:flutter/material.dart';
import 'package:web/models/config/object.dart';
import 'package:web/screens/config/form.dart';
import 'package:web/services/config/config.dart';
import 'package:web/widgets/error.dart';

class ConfigForm extends StatefulWidget {
  final String id;
  final Function() onCancel;
  const ConfigForm({super.key, required this.id, required this.onCancel});

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
    try {
      var res = await ConfigService().getConfig(widget.id);
      setState(() {
        isLoading = false;
        configs = res.data!.configs ?? [];
        error = res.error;
      });
    } on Exception catch (e) {
      print(e);
      setState(() {
        isLoading = false;
        error = "An error occurred. Please try again.";
      });
    }
  }

  void saveConfig() {
    // Save configuration
  }

  @override
  void initState() {
    fetchConfig();
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
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
                    mainAxisAlignment: MainAxisAlignment.end,
                    children: [
                      ElevatedButton(
                        onPressed: () {
                          widget.onCancel();
                        },

                        child: const Text("Cancel"),
                      ),
                      const SizedBox(width: 10),
                      ElevatedButton(
                        onPressed: () {},
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
