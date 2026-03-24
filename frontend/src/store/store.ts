import { configureStore } from "@reduxjs/toolkit";
import { baseApi } from "../api/baseApi";
import { uiSlice } from "./uiSlice";
import { readerSlice } from "./readerSlice";

export const store = configureStore({
  reducer: {
    [baseApi.reducerPath]: baseApi.reducer,
    ui: uiSlice.reducer,
    reader: readerSlice.reducer,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware().concat(baseApi.middleware),
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
