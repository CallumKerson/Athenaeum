/**
 NavigationView.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import SwiftUI
import UniformTypeIdentifiers

struct AppNavigationView: View {
    @Environment(\.importFiles) var importFiles
    @EnvironmentObject private var model: AthenaeumModel

    @ViewBuilder var body: some View {
        SidebarNavigationView()
            .frame(minWidth: 900, maxWidth: .infinity, minHeight: 500, maxHeight: .infinity)
            .toolbar {
                ToolbarItem(placement: .primaryAction) {
                    Button(action: {
                        importFiles.callAsFunction(multipleOfType: [UTType.audio]) { result in
                            switch result {
                            case let .success(urls):
                                urls.printByIndex()
                                for url in urls {
                                    model.addAudiobook(from: url)
                                }
                            case let .failure(error):
                                print(error.localizedDescription)
                            case .none:
                                print("User cancelled")
                            }
                        }
                        print("Import button was tapped")
                    }) {
                        Image(systemName: "plus")
                    }
                }
            }
    }
}

extension Array {
    /// Prints self to std output, with one element per line, prefixed by
    /// the element's index in square brackets.

    public func printByIndex(delimiter: String = " ") {
        for (index, value) in enumerated() {
            print("[\(index)]\(delimiter)\(value)")
        }
    }
}

struct AppNavigationView_Previews: PreviewProvider {
    static var previews: some View {
        AppNavigationView()
    }
}
