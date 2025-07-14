#include <stdio.h>
#include <stdlib.h>
#include <string.h>

// // TODO: maybe this should contain an extra byte to contain the number of letters

const int LETTER_SET_BYTES = 26;
typedef char letter_set[LETTER_SET_BYTES]; // TODO: change this to be 4 bits per letter (max count per letter of 15)

typedef struct dictionary
{
    int num_word_sets;
    letter_set *word_sets;
    // words stored with their anagrams
    char ***anagrams_by_word_set_index;
} dictionary;

void populate_letter_set_from_string(letter_set set, char *string)
{
    for (int i = 0; string[i] != '\0'; i++)
    {
        const int letter_index = string[i] - 'a';
        set[letter_index]++;
    }
}

int letter_set_is_subset(letter_set set1, letter_set set2) // set1 is subset of set2?
{
    for (int i = 0; i < LETTER_SET_BYTES; i++)
    {
        if (set1[i] > set2[i])
        {
            return 0;
        }
    }
    return 1;
}

#define MAX_LINE_LENGTH 128

typedef struct word_summary
{
    int length;
    char *word;
    letter_set letter_set;
} word_summary;

int compare_word_summaries(const void *a, const void *b)
{
    const word_summary *word_summary_a = *(const word_summary **)a;
    const word_summary *word_summary_b = *(const word_summary **)b;
    int length_diff = word_summary_a->length - word_summary_b->length;
    if (length_diff != 0)
    {
        return length_diff;
    }
    return memcmp(word_summary_a->letter_set, word_summary_b->letter_set, LETTER_SET_BYTES);
}

dictionary *create_dictionary(char *filename)
{
    // TODO: this could all be done at compile time

    FILE *file = fopen(filename, "r");
    if (!file)
    {
        perror("Error opening dictionary file");
        return NULL;
    }

    int lines = 0;
    int ch;
    while ((ch = fgetc(file)) != EOF)
    {
        if (ch == '\n')
        {
            lines++;
        }
    }

    rewind(file);

    word_summary *word_summaries = calloc(lines, sizeof(word_summary));
    char line[MAX_LINE_LENGTH];
    int line_index = 0;
    while (fgets(line, sizeof(line), file))
    {
        line[strcspn(line, "\r\n")] = '\0'; // todo make this safe

        word_summary *current_word_summary = &word_summaries[line_index];

        current_word_summary->length = strlen(line);
        populate_letter_set_from_string(current_word_summary->letter_set, line);

        current_word_summary->word = malloc(current_word_summary->length + 1);
        strcpy(current_word_summary->word, line);

        line_index++;
    }

    fclose(file);

    word_summary **word_summaries_pointers = calloc(lines, sizeof(word_summary *));
    for (int i = 0; i < lines; i++)
    {
        word_summaries_pointers[i] = &word_summaries[i];
    }

    qsort(word_summaries_pointers, lines, sizeof(word_summary *), compare_word_summaries);

    dictionary *dictionary = calloc(sizeof(dictionary));
    dictionary->num_word_sets = 0;
    dictionary->word_sets = calloc(lines, sizeof(letter_set));
    dictionary->anagrams_by_word_set_index = calloc(lines, sizeof(char ***));

    for (int i = 0; i < lines; i++)
    {
        word_summary *word_summary = word_summaries_pointers[i];
    }
    return NULL;
}

void usage(const char *progname)
{
    fprintf(stderr, "Usage: %s -dictionary <filename>\n", progname);
    exit(1);
}

int main(int argc, char *argv[])
{
    if (argc != 3 || strcmp(argv[1], "-dictionary") != 0)
    {
        usage(argv[0]);
    }

    const char *filename = argv[2];
    FILE *file = fopen(filename, "r");
    if (!file)
    {
        perror("Error opening dictionary file");
        return 1;
    }

    char line[MAX_LINE_LENGTH];
    while (fgets(line, sizeof(line), file))
    {
        line[strcspn(line, "\r\n")] = '\0';

        // printf("Word: %s\n", line);
    }

    fclose(file);
    return 0;
}