#include <stdio.h>
#include <stdlib.h>
#include <string.h>

// // TODO: maybe this should contain an extra byte to contain the number of letters

#define MIN_WORD_LENGTH 4

const int NUM_LETTERS = 26;
const int letter_distributions[NUM_LETTERS] = {12, 3, 5, 6, 18, 3, 6, 4, 11, 2, 2, 7, 4, 9, 10, 4, 2, 10, 8, 9, 5, 2, 2, 2, 2, 2};

const unsigned int seed = 1234123;

// This code released under creative commons license here:
//
// https://gist.github.com/kevinmoran/0198d8e9de0da7057abe8b8b34d50f86
unsigned int SquirrelNoise5(int positionX, unsigned int seed)
{
    const unsigned int SQ5_BIT_NOISE1 = 0xd2a80a3f; // 11010010101010000000101000111111
    const unsigned int SQ5_BIT_NOISE2 = 0xa884f197; // 10101000100001001111000110010111
    const unsigned int SQ5_BIT_NOISE3 = 0x6C736F4B; // 01101100011100110110111101001011
    const unsigned int SQ5_BIT_NOISE4 = 0xB79F3ABB; // 10110111100111110011101010111011
    const unsigned int SQ5_BIT_NOISE5 = 0x1b56c4f5; // 00011011010101101100010011110101

    unsigned int mangledBits = (unsigned int)positionX;
    mangledBits *= SQ5_BIT_NOISE1;
    mangledBits += seed;
    mangledBits ^= (mangledBits >> 9);
    mangledBits += SQ5_BIT_NOISE2;
    mangledBits ^= (mangledBits >> 11);
    mangledBits *= SQ5_BIT_NOISE3;
    mangledBits ^= (mangledBits >> 13);
    mangledBits += SQ5_BIT_NOISE4;
    mangledBits ^= (mangledBits >> 15);
    mangledBits *= SQ5_BIT_NOISE5;
    mangledBits ^= (mangledBits >> 17);
    return mangledBits;
}

char *choose_letter_flip_order(unsigned int seed, const int *letter_distributions)
{
    int num_letters_to_flip = 0;
    for (int i = 0; i < NUM_LETTERS; i++)
    {
        num_letters_to_flip += letter_distributions[i];
    }

    printf("num_letters_to_flip: %d\n", num_letters_to_flip);
    char *letter_flip_order = calloc(num_letters_to_flip + 1 /* for the null terminator */, sizeof(char));
    int i = 0;
    for (int letter_index = 0; letter_index < NUM_LETTERS; letter_index++)
    {
        const char letter = 'a' + letter_index;
        const int num_letters_to_flip_for_this_letter = letter_distributions[letter_index];
        for (int j = 0; j < num_letters_to_flip_for_this_letter; j++)
        {
            letter_flip_order[i] = letter;
            i++;
        }
    }

    for (int i = 0; i < num_letters_to_flip; i++)
    {
        int random_index = SquirrelNoise5(seed, i) % num_letters_to_flip;
        char tmp = letter_flip_order[i];
        letter_flip_order[i] = letter_flip_order[random_index];
        letter_flip_order[random_index] = tmp;
    }

    return letter_flip_order;
}

const int LETTER_SET_BYTES = 26;
typedef char letter_set[LETTER_SET_BYTES]; // TODO: change this to be 4 bits per letter (max count per letter of 15)

typedef struct word_summary
{
    int length;
    char *word;
    letter_set letter_set;
} word_summary;

typedef struct word_summary_index_range
{
    int start; // inclusive
    int end;   // inclusive
} word_summary_index_range;

typedef struct dictionary
{
    int num_letter_sets;
    letter_set *letter_sets;
    word_summary_index_range *anagram_words_by_letter_set_index;
    int num_words;
    word_summary *words_sorted;
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

int letter_set_is_equal(letter_set set1, letter_set set2)
{
    return memcmp(set1, set2, LETTER_SET_BYTES) == 0;
}

#define MAX_LINE_LENGTH 128

int compare_word_summaries(const void *a, const void *b)
{
    const word_summary *word_summary_a = (const word_summary *)a;
    const word_summary *word_summary_b = (const word_summary *)b;
    int length_diff = word_summary_a->length - word_summary_b->length;
    if (length_diff != 0)
    {
        return -length_diff;
    }
    return memcmp(word_summary_a->letter_set, word_summary_b->letter_set, LETTER_SET_BYTES);
}

dictionary *create_dictionary(const char *filename)
{
    // TODO: this could all be done at compile time

    FILE *file = fopen(filename, "r");
    if (!file)
    {
        perror("Error opening dictionary file");
        return NULL;
    }

    int num_lines = 0;
    int ch;
    while ((ch = fgetc(file)) != EOF)
    {
        if (ch == '\n')
        {
            num_lines++;
        }
    }

    rewind(file);

    word_summary *word_summaries = calloc(num_lines, sizeof(word_summary));
    char line[MAX_LINE_LENGTH];
    int num_words = 0;
    while (fgets(line, sizeof(line), file))
    {
        line[strcspn(line, "\r\n")] = '\0'; // todo make this safe

        const int length = strlen(line);
        if (length < MIN_WORD_LENGTH)
        {
            continue;
        }

        word_summary *current_word_summary = &word_summaries[num_words];

        current_word_summary->length = strlen(line);
        populate_letter_set_from_string(current_word_summary->letter_set, line);

        current_word_summary->word = malloc(current_word_summary->length + 1);
        strcpy(current_word_summary->word, line);

        num_words++;
    }

    fclose(file);

    qsort(word_summaries, num_words, sizeof(word_summary), compare_word_summaries);

    int num_letter_sets_to_allocate = num_words; // TODO: this is extra large, we should allocate fewer
    dictionary *dictionary = calloc(1, sizeof(dictionary));
    dictionary->num_letter_sets = 0;
    dictionary->letter_sets = calloc(num_letter_sets_to_allocate, sizeof(letter_set));
    dictionary->anagram_words_by_letter_set_index = calloc(num_letter_sets_to_allocate, sizeof(word_summary_index_range));
    dictionary->words_sorted = word_summaries;
    dictionary->num_words = num_words;

    for (int i = 0; i < num_words; i++)
    {
        if (i == 0 || !letter_set_is_equal(word_summaries[i].letter_set, word_summaries[i - 1].letter_set))
        {
            int current_letter_set_index = dictionary->num_letter_sets;
            dictionary->num_letter_sets++;

            memcpy(dictionary->letter_sets[current_letter_set_index], word_summaries[i].letter_set, LETTER_SET_BYTES);

            dictionary->anagram_words_by_letter_set_index[current_letter_set_index].start = i;
            dictionary->anagram_words_by_letter_set_index[current_letter_set_index].end = i;
        }
        else
        {
            int current_letter_set_index = dictionary->num_letter_sets - 1;
            dictionary->anagram_words_by_letter_set_index[current_letter_set_index].end = i;
            printf("continuing letter set %s\n", word_summaries[i].word);
        }
    }

    return dictionary;
}

void print_dictionary(dictionary *dictionary)
{
    for (int i = 0; i < dictionary->num_letter_sets; i++)
    {
        for (int j = 0; j < LETTER_SET_BYTES; j++)
        {
            if (j > 0)
            {
                printf("_");
            }
            printf("%d", dictionary->letter_sets[i][j]);
        }
        printf(" ");
        const int start = dictionary->anagram_words_by_letter_set_index[i].start;
        const int end = dictionary->anagram_words_by_letter_set_index[i].end;
        for (int j = start; j <= end; j++)
        {
            printf(" %s ", dictionary->words_sorted[j].word);
        }
        printf("\n");
    }
}

typedef struct player_state
{
    char *name;
    int max_words;
    int num_words;
    int *words;
} player_state;

typedef struct game_state
{
    int num_players;
    player_state *players;
    dictionary *dictionary;
    char *letter_flip_order;
    int num_letters_flipped;
    int turn_number;
    int player_turn_index; // the player who flips the next letter
    // TODO: letters in the center
} game_state;

// TODO: setup game state
// TODO: run game turn

void usage(const char *progname)
{
    fprintf(stderr, "Usage: %s -dictionary <filename>\n", progname);
    fprintf(stderr, "each word in the dictionary should be its own line\n");
    exit(1);
}

int main(int argc, char *argv[])
{
    if (argc != 3 || strcmp(argv[1], "-dictionary") != 0)
    {
        usage(argv[0]);
    }

    // const char *filename = argv[2];

    // const dictionary *dictionary = create_dictionary(filename);
    // if (NULL == dictionary)
    // {
    //     fprintf(stderr, "Error creating dictionary\n");
    //     return 1;
    // }
    // print_dictionary(dictionary);
    char *letter_flip_order = choose_letter_flip_order(141232123, letter_distributions);
    for (int i = 0; letter_flip_order[i] != '\0'; i++)
    {
        printf("%c", letter_flip_order[i]);
    }
    printf("\n");
    return 0;
}