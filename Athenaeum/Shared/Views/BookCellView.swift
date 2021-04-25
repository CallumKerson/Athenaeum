/**
 BookCellView.swift
 Copyright (c) 2021 Callum Kerr-Edwards
 */

import SwiftUI

struct BookCellView: View {
    @ObservedObject var viewModel: BookCellViewModel

    @ViewBuilder
    var body: some View {
        if viewModel.loading {
            ProgressView()
                .onAppear(perform: { viewModel.reload() })
        } else if let error = viewModel.error {
            Label(error.description, systemImage: "exclamationmark.triangle")
        } else if let item = viewModel.book {
            VStack(alignment: .leading) {
                Text(item.title)
                    .font(.headline)
                    .lineLimit(3)
                HStack {
                    Text(item.authorString)
                        .font(.caption)
                        .foregroundColor(.gray)
                }
            }

            .padding(.vertical, 4)
        }
    }
}
