/**
 BookRowView.swift
 Copyright (c) 2020 Callum Kerr-Edwards */

import SwiftUI

struct BookRowView: View {
    var book: Book

    @EnvironmentObject private var model: AthenaeumModel

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

    var verticalRowPadding: CGFloat {
        #if os(macOS)
            return 10
        #else
            return 0
        #endif
    }

    var verticalTextPadding: CGFloat {
        #if os(iOS)
            return 8
        #else
            return 0
        #endif
    }

    var body: some View {
        HStack(alignment: .top) {
            if let coverImage = book.image {
                coverImage
                    .resizable()
                    .aspectRatio(contentMode: .fill)
                    .frame(width: size, height: size)
                    .clipShape(RoundedRectangle(cornerRadius: cornerRadius, style: .continuous))
                    .accessibility(hidden: true)
            }
            VStack(alignment: .leading) {
                Text(book.metadata.title)
                    .font(.headline)
                    .lineLimit(1)
            }
            .padding(.vertical, verticalTextPadding)

            Spacer(minLength: 0)
        }
        .font(.subheadline)
        .padding(.vertical, verticalRowPadding)
        .accessibilityElement(children: .combine)
        .frame(minWidth: 300)
    }
}

struct BookRowView_Previews: PreviewProvider {
    static var previews: some View {
        BookRowView(book: Book(metadata: BookMetadata(title: "Book title")))
    }
}
