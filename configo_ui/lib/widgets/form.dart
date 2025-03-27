import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:configo_ui/models/config/object.dart';
import 'package:configo_ui/models/config/type.dart';
import 'package:configo_ui/widgets/json_editor.dart';

class ConfigFormWidget extends StatefulWidget {
  final List<ConfigObject> configs;

  const ConfigFormWidget({super.key, required this.configs});

  @override
  ConfigFormWidgetState createState() => ConfigFormWidgetState();
}

class ConfigFormWidgetState extends State<ConfigFormWidget> {
  Map<String, dynamic> formValues = {};
  Map<String, bool> isExpandedMap = {};

  @override
  Widget build(BuildContext context) {
    List<ConfigObject> sortedConfigs = _sortConfigs(widget.configs);
    return LayoutBuilder(
      builder: (context, constraints) {
        return Padding(
          padding: const EdgeInsets.only(left: 16.0, bottom: 16.0, top: 16.0),
          child: Wrap(
            spacing: 16.0,
            runSpacing: 16.0,
            children:
                sortedConfigs
                    .map(
                      (config) =>
                          _buildConfigField(config, constraints.maxWidth),
                    )
                    .toList(),
          ),
        );
      },
    );
  }

  List<ConfigObject> _sortConfigs(List<ConfigObject> configs) {
    const typeOrder = {
      ConfigType.string: 1,
      ConfigType.choice: 2,
      ConfigType.number: 3,
      ConfigType.boolean: 4,
      ConfigType.bigText: 5,
      ConfigType.json: 6,
    };
    configs.sort(
      (a, b) => (typeOrder[a.type] ?? 99).compareTo(typeOrder[b.type] ?? 99),
    );
    return configs;
  }

  Widget _buildConfigField(ConfigObject config, double maxWidth) {
    if (config.children != null && config.children!.isNotEmpty) {
      return _buildExpandableConfig(config, maxWidth);
    }
    return _buildInputField(config, maxWidth);
  }

  Widget _buildExpandableConfig(ConfigObject config, double maxWidth) {
    if (!isExpandedMap.containsKey(config.name)) {
      isExpandedMap[config.name] = false;
    }
    return SizedBox(
      width: maxWidth,
      child: Card(
        child: Column(
          children: [
            Container(
              decoration: BoxDecoration(
                border: Border.all(
                  color: Theme.of(context).colorScheme.primary.withOpacity(0.2),
                  width: 1,
                ),
                borderRadius: BorderRadius.circular(12.0),
                color: Theme.of(context).colorScheme.primary.withOpacity(0.1),
              ),
              child: ListTile(
                title: Text(
                  config.name,
                  style: const TextStyle(fontWeight: FontWeight.bold),
                ),
                subtitle: Text(config.description),
                trailing: Icon(
                  isExpandedMap[config.name]!
                      ? Icons.expand_less
                      : Icons.expand_more,
                ),
                onTap: () {
                  setState(() {
                    isExpandedMap[config.name] = !isExpandedMap[config.name]!;
                  });
                },
              ),
            ),
            if (isExpandedMap[config.name]!)
              Padding(
                padding: const EdgeInsets.only(
                  left: 16.0,
                  top: 16.0,
                  bottom: 16.0,
                ),
                child: Wrap(
                  spacing: 16.0,
                  runSpacing: 16.0,
                  children:
                      _sortConfigs(config.children!)
                          .map((child) => _buildConfigField(child, maxWidth))
                          .toList(),
                ),
              ),
          ],
        ),
      ),
    );
  }


  double getFieldWidth(double maxWidth, ConfigType type) {
    switch (type) {
      case ConfigType.bigText:
        return maxWidth > 400 ? 600 : maxWidth * 0.9;
      case ConfigType.json:
        return maxWidth > 400 ? 600 : maxWidth * 0.9;
      case ConfigType.boolean:
        return 300;
      case ConfigType.number:
        return 200; 
      default:
        return maxWidth > 400 ? 400 : maxWidth * 0.9;
    }
  }


  Widget _buildInputField(ConfigObject config, double maxWidth) {
    double fieldWidth = getFieldWidth(maxWidth, config.type);
    return SizedBox(
      width: fieldWidth,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            config.description,
            style: TextStyle(fontSize: 12, color: Colors.grey),
          ),
          const SizedBox(height: 4),
          _getFieldWidget(config),
        ],
      ),
    );
  }

  Widget _getFieldWidget(ConfigObject config) {
    switch (config.type) {
      case ConfigType.string:
        return TextFormField(
          initialValue: config.value ?? '',
          decoration: InputDecoration(
            labelText: config.name,
            border: OutlineInputBorder(),
          ),
          onChanged: (value) => config.value = value,
        );
      case ConfigType.bigText:
        return TextFormField(
          initialValue: config.value ?? '',
          decoration: InputDecoration(
            labelText: config.name,
            border: OutlineInputBorder(),
          ),
          minLines: 3,
          maxLines: 5,
          onChanged: (value) => config.value = value,
        );
      case ConfigType.number:
        return TextFormField(
          keyboardType: const TextInputType.numberWithOptions(decimal: true),
          initialValue: config.value == null ? '' : config.value.toString(),
          decoration: InputDecoration(
            labelText: config.name,
            border: OutlineInputBorder(),
          ),
          inputFormatters: [
            FilteringTextInputFormatter.allow(
              RegExp(r'^\d*\.?\d*$'),
            ), // Only allows numbers and a single decimal point
          ],
          onChanged: (value) {
            //convert to float
            config.value = value.isEmpty ? null : double.parse(value);
          },
        );

      case ConfigType.boolean:
        return SwitchListTile(
          title: Text(config.name),
          value: config.value,
          onChanged: (value) {
            setState(() {
              config.value = value;    
            });
          },
        );
      case ConfigType.choice:
        return DropdownButtonFormField<String>(
          value: config.value ?? (config.value == "" ? null : ''),
          decoration: InputDecoration(
            labelText: config.name,
            border: OutlineInputBorder(),
          ),
          items:
              config.choices!
                  .map(
                    (choice) =>
                        DropdownMenuItem(value: choice, child: Text(choice)),
                  )
                  .toList(),
          onChanged: (value) => config.value = value,
        );
      case ConfigType.json:
        return JsonEditorWidget(
          name: config.name,
          initialValue: config.value != null ? config.value.toString() : '',
          onChanged: (value) => {config.value = value},
        );
      default:
        return const SizedBox.shrink();
    }
  }
}
