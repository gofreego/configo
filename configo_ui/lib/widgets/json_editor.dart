import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:flutter_code_editor/flutter_code_editor.dart';
import 'package:flutter_highlight/themes/a11y-dark.dart';
import 'package:highlight/languages/json.dart';

class JsonEditorWidget extends StatefulWidget {
  final String? initialValue;
  final String? name;
  final int minLines;
  final int maxLines;
  final ValueChanged<String>? onChanged;

  const JsonEditorWidget({
    super.key,
    this.initialValue,
    this.name,
    this.minLines = 5,
    this.maxLines = 5,
    this.onChanged,
  });

  @override
  State<JsonEditorWidget> createState() => _JsonEditorWidgetState();
}

class _JsonEditorWidgetState extends State<JsonEditorWidget> {
  late CodeController _controller;

  @override
  void initState() {
    super.initState();
    _controller = CodeController(
      text: widget.initialValue ?? '', // Use initial value if provided
      language: json, // Set language as JSON
    );

    _controller.addListener(() {
      if (widget.onChanged != null) {
        widget.onChanged!(_controller.text);
      }
    });
  }

  void _formatJson() {
    if (_controller.text.trim().isEmpty) return; // Skip validation if empty
    try {
      final parsedJson = jsonDecode(_controller.text);
      final formattedJson = const JsonEncoder.withIndent(
        '  ',
      ).convert(parsedJson);
      setState(() {
        _controller.text = formattedJson;
      });
      if (widget.onChanged != null) {
        widget.onChanged!(formattedJson);
      }
    } catch (e) {
      var msg = e.toString();
      ScaffoldMessenger.of(
        context,
      ).showSnackBar(SnackBar(content: Text("Invalid JSON: $msg")));
    }
  }

  String getJsonValue() {
    return _controller.text;
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        if (widget.name != null)
          Padding(
            padding: const EdgeInsets.symmetric(vertical: 4.0),
            child: Text(
              widget.name!,
              style: const TextStyle(
                fontSize: 16,
                // fontWeight: FontWeight.bold,
              ),
            ),
          ),
        CodeTheme(
          data: CodeThemeData(styles: a11yDarkTheme),
          child: SizedBox(
            height: 26.0 * widget.maxLines, // Set height based on maxLines
            child: Stack(
              children: [
                SingleChildScrollView(
                  child: CodeField(
                    controller: _controller,
                    minLines:
                        widget.minLines, // Minimum number of lines to show
                  ),
                ),
                Align(
                  alignment: Alignment.topRight,
                  child: IconButton(  
                    icon: const Icon(Icons.format_align_left),
                    onPressed: _formatJson,
                    tooltip: "Format JSON",
                  ),
                ),
              ],
            ),
          ),
        ),
      ],
    );
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }
}
