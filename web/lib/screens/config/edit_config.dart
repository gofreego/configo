import 'package:flutter/material.dart';
import 'package:web/screens/config/form.dart';

class ConfigForm extends StatefulWidget {
  final String id;
  final Function() onCancel;
  const ConfigForm({super.key, required this.id, required this.onCancel});
  

  @override
  State<ConfigForm> createState() => _ConfigFormState();
}

class _ConfigFormState extends State<ConfigForm> {
  bool isLoading = true;

  void fetchConfig() {
    setState(() {
      isLoading = true;
    });

    // Simulate network delay
    Future.delayed(const Duration(seconds: 1), () {
      setState(() {
        isLoading = false;
      });
    });
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
              : Column(
                children: [
                  NewWidget(),
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
                        child: const Text("Submit"),
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
