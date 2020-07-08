/**
 BookDetailView.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import SwiftUI

struct BookDetailView: View {
    
    var book: Book
    
    var size: CGFloat {
        #if os(iOS)
            return 96
        #else
            return 60
        #endif
    }

    var cornerRadius: CGFloat {
        #if os(iOS)
            return 16
        #else
            return 8
        #endif
    }
    
    var body: some View {
        VStack(alignment: HorizontalAlignment.leading, spacing: 12) {
            HStack(alignment: VerticalAlignment.center, spacing: 24) {
                if let coverImage = book.image {
                    coverImage
                        .resizable()
                        .aspectRatio(contentMode: .fill)
                        .frame(width: size, height: size)
                        .clipShape(RoundedRectangle(cornerRadius: cornerRadius, style: .continuous))
                        .accessibility(hidden: true)
                }
                VStack(alignment: HorizontalAlignment.leading) {
                    BookTextView(metadata: book.metadata)
                    Spacer()
                }
            }
        }
    }
}

struct BookDetailView_Previews: PreviewProvider {
    static var previews: some View {
        BookDetailView(book: Book(metadata: BookMetadata(title: "Book title")))
    }
}
