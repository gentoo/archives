$("#show-more-replies").click(function () {
    $(".more-replies").removeClass("d-none");
    $("#show-more-replies").addClass("d-none");
    $("#show-less-replies").removeClass("d-none");
});

$("#show-less-replies").click(function () {
    $(".more-replies").addClass("d-none");
    $("#show-less-replies").addClass("d-none");
    $("#show-more-replies").removeClass("d-none");
});
