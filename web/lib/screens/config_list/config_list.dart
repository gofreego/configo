import 'package:flutter/material.dart';
import 'package:web/models/config/metadata.dart';
import 'package:web/screens/config_list/service_info.dart';
import 'package:web/services/config/config.dart';

class ListConfigScreen extends StatefulWidget {
  const ListConfigScreen({super.key});

  @override
  State<ListConfigScreen> createState() => _ListConfigScreenState();
}

class _ListConfigScreenState extends State<ListConfigScreen> {
  bool isLoading = true;
  String? error;
  ConfigMetadataResponse? metadata;

  @override
  void initState() {
    ConfigService()
        .getConfigMetadata()
        .then((value) {
          setState(() {
            isLoading = false;
            metadata = value;
          });
        })
        .catchError((error) {
          setState(() {
            isLoading = false;
            this.error = error.toString();
          });
        });
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return ServiceScreen();
  }
}

class ServiceScreen extends StatefulWidget {
  const ServiceScreen({super.key});

  @override
  State<ServiceScreen> createState() => _ServiceScreenState();
}

class _ServiceScreenState extends State<ServiceScreen> {
  final String serviceName = "Test Service";
  final String serviceDescription =
      "A powerful and scalable testing solution for APIs, applications, and security assessments. It provides end-to-end testing capabilities to ensure the reliability, security, and performance of your applications.";

  final List<Map<String, String>> serviceItems = [
    {
      "title": "Automated Testing",
      "description": "Runs test cases automatically with CI/CD integration.",
    },
    {
      "title": "API Testing",
      "description":
          "Validates REST, GraphQL, and gRPC services for reliability.",
    },
    {
      "title": "Performance Testing",
      "description": "Ensures system stability under high loads.",
    },
    {
      "title": "Security Testing",
      "description": "Identifies vulnerabilities and security threats.",
    },
    {
      "title": "Exam Assessments",
      "description": "Conducts and evaluates online exams.",
    },
  ];

  Map<int, bool> expandedStates = {};
  Map<int, bool> isLoading = {};

  void fetchConfig(int index) {
    setState(() {
      isLoading[index] = true;
      expandedStates[index] = true;
    });

    // Simulate network delay
    Future.delayed(const Duration(seconds: 1), () {
      setState(() {
        isLoading[index] = false;
        expandedStates[index] = true;
      });
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            ServiceInfoWidget(),

            const SizedBox(height: 20),
            Expanded(
              child: ListView.builder(
                itemCount: serviceItems.length,
                itemBuilder: (context, index) {
                  final item = serviceItems[index];
                  final isExpanded = expandedStates[index] ?? false;
                  final loading = isLoading[index] ?? false;
                  return Card(
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(12),
                    ),
                    elevation: 3,
                    margin: const EdgeInsets.symmetric(vertical: 8),
                    child: Column(
                      children: [
                        InkWell(
                          onTap: () {
                            if (!isExpanded) {
                              fetchConfig(index);
                            } else {
                              setState(() {
                                expandedStates[index] = false;
                              });
                            }
                          },
                          borderRadius: BorderRadius.circular(12),
                          child: Padding(
                            padding: const EdgeInsets.symmetric(
                              vertical: 12,
                              horizontal: 16,
                            ),
                            child: Row(
                              mainAxisAlignment:
                                  MainAxisAlignment.spaceBetween,
                              children: [
                                Row(
                                  children: [
                                     Icon(
                                      Icons.settings,
                                      color: Theme.of(context).colorScheme.primary,
                                    ),
                                    const SizedBox(width: 12),
                                    Text(
                                      item["title"]!,
                                      style: const TextStyle(
                                        fontWeight: FontWeight.bold,
                                        fontSize: 18,
                                      ),
                                    ),
                                  ],
                                ),
                                Icon(
                                  isExpanded
                                      ? Icons.expand_less
                                      : Icons.expand_more,
                                ),
                              ],
                            ),
                          ),
                        ),
                        AnimatedSize(
                          duration: const Duration(milliseconds: 200),
                          curve: Curves.easeInOut,
                          child:
                              isExpanded
                                  ? Padding(
                                    padding: const EdgeInsets.symmetric(
                                      horizontal: 16,
                                      vertical: 8,
                                    ),
                                    child: Column(
                                      children: [
                                        loading
                                            ? const Center(
                                              child: Padding(
                                                padding: EdgeInsets.all(8.0),
                                                child:
                                                    CircularProgressIndicator(),
                                              ),
                                            )
                                            : Column(
                                              children: [
                                                TextField(
                                                  decoration: InputDecoration(
                                                    labelText:
                                                        "Enter Configuration",
                                                    border: OutlineInputBorder(
                                                      borderRadius:
                                                          BorderRadius.circular(
                                                            8,
                                                          ),
                                                    ),
                                                  ),
                                                ),
                                                const SizedBox(height: 10),
                                                ElevatedButton(
                                                  onPressed: () {},
                                                  child: const Text("Submit"),
                                                ),
                                              ],
                                            ),
                                      ],
                                    ),
                                  )
                                  : const SizedBox.shrink(),
                        ),
                      ],
                    ),
                  );
                },
              ),
            ),
          ],
        ),
      ),
    );
  }
}
