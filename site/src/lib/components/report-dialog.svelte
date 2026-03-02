<script>
    const {url} = $props();
    const ao3Id = url.split('/').slice(-1).join('-');
</script>

<button command="show-modal" commandfor="report-dialog-{ao3Id}">Open dialog</button>

<dialog id="report-dialog-{ao3Id}" closedby="any">
    <form method="POST" action="?/report">
        <input type="hidden" name="bookUrl" value={url} required/>
        <select name="reason" id="reason" required >
            <option value="unavailable">This story is no longer online</option>
            <option value="unacceptable">This isn't a story</option>
        </select>
        <input type="text" name="comment" required/>
        <input type="submit" value="report"/>

        <button commandfor="report-dialog-{ao3Id}" command="close">Close</button>
    </form>
</dialog>
