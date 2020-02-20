/**
 LibraryView.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import RealmSwift
import SwiftUI

struct LibraryView: View {
    @ObservedObject var library: Library

    init(withLibrary library: Library = Library.global) {
        self.library = library
    }

    var body: some View {
        NavigationView {
            List {
                Section(header: HeaderView()) {
                    ForEach(library.🎧📚.sorted(by: { $0.author < $1.author })) { book in
                        NavigationLink(destination: AudiobookView(book: book)) {
                            AudiobookCellView(book: book)
                        }
                    }
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
        LibraryView()
    }
}
