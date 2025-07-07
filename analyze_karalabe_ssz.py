#!/usr/bin/env python3
import os
import re
from pathlib import Path

def find_structs_with_ssz(file_path):
    """Find all structs in a file that have SSZ methods."""
    with open(file_path, 'r') as f:
        content = f.read()
    
    # Check if file imports karalabe/ssz
    if 'github.com/karalabe/ssz' not in content:
        return []
    
    # Find all struct definitions
    struct_pattern = r'^type\s+(\w+)\s+struct\s*{' 
    structs = re.findall(struct_pattern, content, re.MULTILINE)
    
    results = []
    for struct_name in structs:
        # Check if this struct has SSZ methods
        ssz_methods = []
        
        # Look for SSZ methods
        method_patterns = [
            (r'func\s*\([^)]*' + struct_name + r'[^)]*\)\s*DefineSSZ\s*\(', 'DefineSSZ'),
            (r'func\s*\([^)]*' + struct_name + r'[^)]*\)\s*SizeSSZ\s*\(', 'SizeSSZ'),
            (r'func\s*\([^)]*' + struct_name + r'[^)]*\)\s*HashTreeRoot\s*\(', 'HashTreeRoot'),
            (r'func\s*\([^)]*' + struct_name + r'[^)]*\)\s*MarshalSSZ\s*\(', 'MarshalSSZ'),
            (r'func\s*\([^)]*' + struct_name + r'[^)]*\)\s*ValidateAfterDecodingSSZ\s*\(', 'ValidateAfterDecodingSSZ'),
        ]
        
        for pattern, method_name in method_patterns:
            if re.search(pattern, content):
                ssz_methods.append(method_name)
        
        # Check if implements ssz.StaticObject
        implements_static = bool(re.search(r'ssz\.StaticObject.*=.*' + struct_name, content))
        
        if ssz_methods:
            results.append({
                'struct': struct_name,
                'methods': ssz_methods,
                'implements_static': implements_static,
                'file': file_path
            })
    
    return results

def main():
    base_path = '/Users/fridrikasmundsson/workspace/beacon-kit-karalabe'
    all_results = []
    
    # Walk through all Go files
    for root, dirs, files in os.walk(base_path):
        for file in files:
            if file.endswith('.go') and not file.endswith('_test.go'):
                file_path = os.path.join(root, file)
                results = find_structs_with_ssz(file_path)
                all_results.extend(results)
    
    # Sort by file path for better organization
    all_results.sort(key=lambda x: x['file'])
    
    print("# Structs using karalabe/ssz for serialization")
    print("=" * 80)
    print()
    
    # Group by directory
    current_dir = None
    for result in all_results:
        dir_name = os.path.dirname(result['file']).replace(base_path, '')
        if dir_name != current_dir:
            current_dir = dir_name
            print(f"\n## Directory: {dir_name or '/'}")
            print("-" * 40)
        
        print(f"\n### {result['struct']}")
        print(f"**File:** `{os.path.basename(result['file'])}`")
        if result['implements_static']:
            print("**Implements:** `ssz.StaticObject`")
        print("**Methods:**")
        for method in result['methods']:
            print(f"  - {method}")
    
    print(f"\n\n## Summary")
    print(f"Total structs using karalabe/ssz: {len(all_results)}")
    
    # Count by directory
    dir_counts = {}
    for result in all_results:
        dir_name = os.path.dirname(result['file']).replace(base_path, '') or '/'
        dir_counts[dir_name] = dir_counts.get(dir_name, 0) + 1
    
    print("\nBreakdown by directory:")
    for dir_name, count in sorted(dir_counts.items()):
        print(f"  - {dir_name}: {count} structs")

if __name__ == '__main__':
    main()