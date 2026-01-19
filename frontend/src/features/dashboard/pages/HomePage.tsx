import { useMemo } from "react";
import { useAuth } from "../../auth/context/AuthContext";
import { useChats } from "../../chats/hooks/useChats";
import { useContacts } from "../../contacts/hooks/useContacts";

export function HomePage() {
  const { logout } = useAuth();
  const { data: chats, loading: chatsLoading } = useChats();
  const { data: contacts, loading: contactsLoading } = useContacts();

  const statusLabel = useMemo(
    () => ({
      online: "Online",
      away: "Away",
      offline: "Offline",
    }),
    []
  );

  function getInitials(name: string) {
    const parts = name.trim().split(/\s+/);
    const first = parts[0]?.[0] ?? "";
    const second = parts[1]?.[0] ?? "";
    return `${first}${second}`.toUpperCase();
  }

  return (
    <div className="container">
      <div className="dashboard">
        <section className="panel">
          <header className="panel-header">
            <div>
              <h1 className="panel-title">Chats</h1>
              <p className="muted">Recent conversations</p>
            </div>
            <button className="secondary" onClick={logout}>
              Logout
            </button>
          </header>

          {chatsLoading ? (
            <div className="empty">Loading chats...</div>
          ) : chats.length === 0 ? (
            <div className="empty">No chats yet.</div>
          ) : (
            <ul className="list" role="list">
              {chats.map((chat) => (
                <li key={chat.id} className="list-item">
                  <div className="avatar">{getInitials(chat.title)}</div>
                  <div className="list-main">
                    <div className="list-row">
                      <span className="list-title">{chat.title}</span>
                      <span className="muted small">{chat.updatedAt}</span>
                    </div>
                    <div className="list-row">
                      <span className="list-sub">{chat.lastMessage}</span>
                      {chat.unreadCount > 0 && (
                        <span className="badge">{chat.unreadCount}</span>
                      )}
                    </div>
                  </div>
                </li>
              ))}
            </ul>
          )}
        </section>

        <section className="panel">
          <header className="panel-header">
            <div>
              <h2 className="panel-title">Contacts</h2>
              <p className="muted">People you talk with</p>
            </div>
          </header>

          {contactsLoading ? (
            <div className="empty">Loading contacts...</div>
          ) : contacts.length === 0 ? (
            <div className="empty">No contacts yet.</div>
          ) : (
            <ul className="list" role="list">
              {contacts.map((contact) => (
                <li key={contact.id} className="list-item">
                  <div className="avatar">{getInitials(contact.name)}</div>
                  <div className="list-main">
                    <div className="list-row">
                      <span className="list-title">{contact.name}</span>
                      <span className={`status-dot status-${contact.status}`} />
                    </div>
                    <div className="list-row">
                      <span className="list-sub">
                        {contact.status === "offline" && contact.lastSeen
                          ? `Last seen ${contact.lastSeen}`
                          : statusLabel[contact.status]}
                      </span>
                    </div>
                  </div>
                </li>
              ))}
            </ul>
          )}
        </section>
      </div>
    </div>
  );
}
