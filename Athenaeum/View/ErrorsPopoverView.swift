/**
 ErrorsPopoverView.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import SwiftUI

struct ErrorsPopoverView: View {
    @ObservedObject var viewModel: ErrorPopoverViewModel

    init(viewModel: ErrorPopoverViewModel) {
        self.viewModel = viewModel
    }

    var body: some View {
        VStack {
            Text("Errors occured:")
                .padding(.bottom)

            ForEach(viewModel.errors) { error in
                VStack {
                    Text(error.path)
                    Text(error.message)
                }.padding(.bottom)
            }

            Button(action: {
                self.viewModel.clearErrors()
            }) {
                Text("Clear")
            }
        }.padding()
    }
}

#if DEBUG
    struct ErrorsPopoverView_Previews: PreviewProvider {
        static var previews: some View {
            ErrorsPopoverView(viewModel: ErrorPopoverViewModel(store: sampleStore))
        }
    }
#endif
