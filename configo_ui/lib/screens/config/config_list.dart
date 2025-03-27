import 'package:flutter/material.dart';
import 'package:configo_ui/models/config/metadata.dart';
import 'package:configo_ui/screens/config/edit_config.dart';
import 'package:configo_ui/widgets/service_info.dart';
import 'package:configo_ui/services/config/config.dart';
import 'package:configo_ui/widgets/error.dart';

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
    var res = await ConfigService().getConfigMetadata();
    setState(() {
      isLoading = false;
      error = res.error;
      metadata = res.data;
    });
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
                ? CustomErrorWidget(
                  errorMessage: error!,
                  onRetry: fetchConfigMetadata,
                )
                : Column(
                  children: [
                    Expanded(
                      child: SingleChildScrollView(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            ServiceInfoWidget(
                              serviceInfo: metadata!.serviceInfo,
                            ),
                            const SizedBox(height: 20),
                            ListView.builder(
                              shrinkWrap: true,
                              physics: NeverScrollableScrollPhysics(),
                              itemCount: metadata!.configInfo.configKeys.length,
                              itemBuilder: (context, index) {
                                final item =
                                    metadata!.configInfo.configKeys[index];
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
      // shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      elevation: 3,
      margin: const EdgeInsets.symmetric(vertical: 8),
      child: Container(
        decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(12),
          border: Border.all(
            color: Theme.of(context).colorScheme.primary.withOpacity(0.5),
            width: 2,
          ),
          // color: Theme.of(context).colorScheme.primary.withOpacity(0.1),
        ),
        child: Column(
          children: [
            InkWell(
              onTap: () {
                setState(() {
                  isExpanded = !isExpanded;
                });
              },
              child: Container(
                decoration: BoxDecoration(
                  borderRadius:
                      isExpanded
                          ? BorderRadius.only(
                            topLeft: Radius.circular(12),
                            topRight: Radius.circular(12),
                          )
                          : BorderRadius.circular(12),
                  color: Theme.of(context).colorScheme.primary.withOpacity(0.1),
                ),
                child: Padding(
                  padding: const EdgeInsets.symmetric(
                    vertical: 12,
                    horizontal: 16,
                  ),
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
            ),
            AnimatedSize(
              duration: const Duration(milliseconds: 300),
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
      ),
    );
  }
}
