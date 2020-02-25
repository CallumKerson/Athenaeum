/**
 LibraryView.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import SwiftUI

struct LibraryView<Lib>: View where Lib: Library {
    @ObservedObject var library: Lib

    init(withLibrary library: Lib) {
        self.library = library
    }

    var body: some View {
        NavigationView {
            List {
                Section(header: HeaderView()) {
                    ForEachAudiobookView(inLibrary: library)
                }
            }
            .frame(minWidth: 425, maxWidth: 425)
        }
        .listStyle(SidebarListStyle())
        .frame(minWidth: 850,
               maxWidth: 850,
               minHeight: 400,
               maxHeight: .infinity)
    }
}

struct HeaderView: View {
    var body: some View {
        VStack {
            HStack(spacing: 20) {
                Text("Library")
                    .layoutPriority(1)
                    .font(.largeTitle)
                Spacer()
            }
        }.padding(.bottom)
    }
}

struct LibraryView_Previews: PreviewProvider {
    static var previews: some View {
        Group {
            LibraryView(withLibrary: MockLibrary())
                .environment(\.colorScheme, .light)
            LibraryView(withLibrary: MockLibrary())
                .environment(\.colorScheme, .dark)
        }
    }
}
