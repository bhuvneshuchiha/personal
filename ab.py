from PyPDF2 import PdfWriter, PdfReader

def split_pdf_to_two(filename,page_number):
    pdf_reader = PdfReader(open(filename, "rb"))
    try:
        assert page_number < len(pdf_reader.pages)
        pdf_writer1 = PdfWriter()
        pdf_writer2 = PdfWriter()

        for page in range(page_number):
            pdf_writer1.add_page(pdf_reader.pages[page])

        for page in range(page_number,len(pdf_reader.pages)):
            pdf_writer2.add_page(pdf_reader.pages[page])

        with open("part1.pdf", 'wb') as file1:
            pdf_writer1.write(file1)

        with open("part2.pdf", 'wb') as file2:
            pdf_writer2.write(file2)

    except AssertionError as e:
        print("Error: The PDF you are cutting has less pages than you want to cut!")

split_pdf_to_two("epod.pdf", 1)
