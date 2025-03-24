import 'package:flutter/material.dart';
import 'package:web/models/config/metadata.dart';
import 'package:web/screens/config/edit_config.dart';
import 'package:web/widgets/service_info.dart';
import 'package:web/services/config/config.dart';
import 'package:web/widgets/error.dart';

class ListConfigScreen extends StatefulWidget {
  const ListConfigScreen({super.key});

  @override
  State<ListConfigScreen> createState() => _ListConfigScreenState();
}

class _ListConfigScreenState extends State<ListConfigScreen> {
  bool isLoading = true;
  String? error;
  ConfigMetadataResponse? metadata;

  void fetchConfigMetadata() async {
    setState(() {
      isLoading = true;
    });
    try {
      var res = await ConfigService().getConfigMetadata();
      setState(() {
        isLoading = false;
        metadata = res.data;
        error = res.error;
      });
    } on Exception catch (_) {
      setState(() {
        isLoading = false;
        error = "An error occurred. Please try again later.";
      });
    }
  }

  @override
  void initState() {
    fetchConfigMetadata();
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child:
            isLoading
                ? const Center(child: CircularProgressIndicator())
                : error != null
                ? CustomErrorWidget(errorMessage: error!, onRetry: fetchConfigMetadata)
                : Column(
                  children: [
                    Expanded(
                      child: SingleChildScrollView(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            ServiceInfoWidget(serviceInfo: metadata!.serviceInfo),
                            const SizedBox(height: 20),
                            ListView.builder(
                              shrinkWrap: true,
                              physics: NeverScrollableScrollPhysics(),
                              itemCount: metadata!.configInfo.configKeys.length,
                              itemBuilder: (context, index) {
                                final item = metadata!.configInfo.configKeys[index];
                                return ConfigTile(id: item);
                              },
                            ),
                          ],
                        ),
                      ),
                    ),
                  ],
                ),
      ),
    );
  }
}

class ConfigTile extends StatefulWidget {
  final String id;
  const ConfigTile({super.key, required this.id});

  @override
  State<ConfigTile> createState() => _ConfigTileState();
}

class _ConfigTileState extends State<ConfigTile> {
  bool isExpanded = false;

  @override
  Widget build(BuildContext context) {
    return Card(
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      elevation: 3,
      margin: const EdgeInsets.symmetric(vertical: 8),
      child: Column(
        children: [
          InkWell(
            onTap: () {
              setState(() {
                isExpanded = !isExpanded;
              });
            },
            borderRadius: BorderRadius.circular(12),
            child: Padding(
              padding: const EdgeInsets.symmetric(vertical: 12, horizontal: 16),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Row(
                    children: [
                      Icon(
                        Icons.settings,
                        color: Theme.of(context).colorScheme.primary,
                      ),
                      const SizedBox(width: 12),
                      Text(
                        widget.id,
                        style: const TextStyle(
                          fontWeight: FontWeight.bold,
                          fontSize: 18,
                        ),
                      ),
                    ],
                  ),
                  Icon(isExpanded ? Icons.expand_less : Icons.expand_more),
                ],
              ),
            ),
          ),
          AnimatedSize(
            duration: const Duration(milliseconds: 200),
            curve: Curves.easeInOut,
            child:
                isExpanded
                    ? ConfigForm(
                      id: widget.id,
                      onCancel: () {
                        setState(() {
                          isExpanded = false;
                        });
                      },
                    )
                    : const SizedBox.shrink(),
          ),
        ],
      ),
    );
  }
}
