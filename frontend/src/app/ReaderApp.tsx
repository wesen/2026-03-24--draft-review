import { useParams } from "react-router-dom";
import { MacWindow, MenuBar } from "../chrome";
import { ReaderPage } from "../reader";
import { useGetReaderLinkQuery } from "../api/readerApi";

export function ReaderApp() {
  const { token } = useParams<{ token: string }>();
  const { data, isLoading, error } = useGetReaderLinkQuery(token || "", {
    skip: !token,
  });

  const menus = [
    { label: "View", items: [{ label: "Reading Mode" }] },
    { label: "Help", items: [{ label: "How to Leave Feedback" }] },
  ];

  if (isLoading) {
    return (
      <div className="dr-desktop">
        <MenuBar menus={menus} />
        <MacWindow title="Loading..." maximized>
          <div
            style={{
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              height: "100%",
              fontFamily: "var(--dr-font-body)",
              fontSize: "var(--dr-font-size-md)",
            }}
          >
            Loading your review...
          </div>
        </MacWindow>
      </div>
    );
  }

  if (error || !data) {
    return (
      <div className="dr-desktop">
        <MenuBar menus={menus} />
        <MacWindow title="Error" maximized>
          <div
            style={{
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              height: "100%",
              fontFamily: "var(--dr-font-body)",
              fontSize: "var(--dr-font-size-md)",
              flexDirection: "column",
              gap: 8,
            }}
          >
            <div style={{ fontSize: 32 }}>{"\u26A0"}</div>
            <div>This review link is invalid or has expired.</div>
          </div>
        </MacWindow>
      </div>
    );
  }

  return (
    <div className="dr-desktop">
      <MenuBar
        menus={menus}
        rightStatus={`Reading as ${data.reader.name}`}
      />
      <MacWindow
        title={`${data.article.title} \u2014 Reader View`}
        maximized
      >
        <ReaderPage article={data.article} />
      </MacWindow>
    </div>
  );
}
