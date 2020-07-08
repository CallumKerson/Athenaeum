/**
 BookDetailView.swift
 Copyright (c) 2020 Callum Kerr-Edwards */

import SwiftUI

struct BookDetailView: View {
    var book: Book

    var cornerRadius: CGFloat {
        #if os(iOS)
            return 16
        #else
            return 8
        #endif
    }

    var body: some View {
        VStack(alignment: HorizontalAlignment.leading, spacing: 12) {
            HStack(alignment: VerticalAlignment.top, spacing: 24) {
                book.image
                    .resizable()
                    .frame(idealWidth: 300, maxWidth: 400, idealHeight: 300, maxHeight: 400)
                    .aspectRatio(contentMode: .fill)
                    .clipShape(RoundedRectangle(cornerRadius: cornerRadius, style: .continuous))
                    .accessibility(hidden: true)

                VStack(alignment: HorizontalAlignment.leading) {
                    BookTextView(metadata: book.metadata)
                        .frame(minWidth: 100)
                    Spacer()
                }
            }
        }.padding()
            .toolbar {
                ToolbarItem {
                    Button("Edit") {
                        print("Edited")
                    }
                }
            }
    }
}

struct BookDetailView_Previews: PreviewProvider {
    static var previews: some View {
        var sampleMetadata = BookMetadata(title: "Murder on the Orient Express")
        sampleMetadata.authors = ["Agatha Christie"]
        sampleMetadata.narrators = ["David Suchet"]
        return BookDetailView(book: Book(metadata: sampleMetadata))
    }
}
