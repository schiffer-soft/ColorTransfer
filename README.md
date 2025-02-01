# Fusion 360 to PrusaSlicer 3MF Converter

![ColorTransfer v0 1 0 32](https://github.com/user-attachments/assets/9b25bbde-6d54-413d-8702-053ee55e70e7)

## Overview
This tool converts Fusion 360 exports (in 3MF format) into a PrusaSlicer 3MF project file with assigned extruders. The tool supports up to 5 different colors.

## Features
- Converts Fusion 360 3MF files to PrusaSlicer-compatible 3MF project files.
- Automatically assigns extruders based on color or other criteria.
- Supports 5 different colors for multi-material printing.
- Allows sorting of colors before conversion.

## Usage
1. Launch `Fusion360`. Open the project you want to export. Select and isolate the individual parts that you want to export. Please note that the exported element may <ins>only have 5 colors<ins>.
    
![grafik](https://github.com/user-attachments/assets/077c86b6-4f5c-4822-b944-6a2f5bf85320)

2. Now right-click on the appropriate component in the project browser and then click on `Save as mesh` in the context menu that appears

![grafik](https://github.com/user-attachments/assets/f775aea9-2952-4f62-beb5-13e161138c24)

3. Now set the format to `3mf` in the menu on the right and click `OK` to save the file.

![grafik](https://github.com/user-attachments/assets/e03f913c-cce1-4abb-a321-a3536f274ae5)

4. Now launch the `ColorTransferGui`. Use the open button `[...]` to open the previously exported `3mf` file.
   If the selected file does not contain any errors, the colors transferred by `Fusion360` are displayed in the left column. On the right-hand side, you can now adjust the order of the extruders using the `[˄]` and `[˅]` buttons. Then save the new file to the desired location using the `Convert and Export` button.  
   
![grafik](https://github.com/user-attachments/assets/d6b91987-56e6-4726-b9d6-ee94119e90c4)

5. As a final step, you can now open the exported file in the `Prusa Slicer`. The extruders are now sorted according to your specifications and the colors are mapped.

![grafik](https://github.com/user-attachments/assets/8a1529af-c2bb-4363-9985-3d5b5b2cdfd5)

## Limitations
- Only 5 colors (extruders) are supported.
- Requires Fusion 360 to export the initial 3MF file.
- PrusaSlicer must be configured to support multi-material printing.
