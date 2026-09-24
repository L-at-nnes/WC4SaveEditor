using WC4SaveEditor.Core.Models;
using WC4SaveEditor.Core.SaveFile;

namespace WC4SaveEditor.Core.Operations;

public sealed class ChangeQueue
{
    private readonly List<PendingChange> _pending = [];

    public IReadOnlyList<PendingChange> Pending => _pending;
    public int Count => _pending.Count;

    public void Add(PendingChange change) => _pending.Add(change);

    public void RemoveAll<T>() where T : PendingChange => _pending.RemoveAll(c => c is T);

    public void RemoveAll(Func<PendingChange, bool> predicate) => _pending.RemoveAll(c => predicate(c));

    public void Clear() => _pending.Clear();

    /// <summary>
    /// Applies every pending change and writes the result to disk. Runs a dry run against a
    /// throwaway clone first: if any change would throw, it throws here before the real
    /// document or the file on disk are touched at all, so a bad queued change can never
    /// leave the document half-mutated while still being listed as pending.
    /// </summary>
    public void Commit(SaveDocument doc)
    {
        ApplyPending(doc.Clone());

        ApplyPending(doc);
        SaveFileWriter.Commit(doc);
        _pending.Clear();
    }

    /// <summary>
    /// Applies every pending change to <paramref name="doc"/> without writing to disk or
    /// clearing the queue - used to render a live "what the map will look like" preview.
    /// Callers must pass a throwaway document (e.g. freshly re-read from disk), never the
    /// document the rest of the app is editing, since this mutates it in place.
    /// </summary>
    public void ApplyToPreview(SaveDocument doc) => ApplyPending(doc);

    private void ApplyPending(SaveDocument doc)
    {
        foreach (var change in _pending.Where(c => !c.IsStructural))
        {
            change.Apply(doc);
        }
        foreach (var change in _pending.Where(c => c.IsStructural))
        {
            change.Apply(doc);
        }
    }
}
