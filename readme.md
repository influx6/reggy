#Reggy
Reggy is a package for creating simple url pattern matchers using two style and provides validation methods that can extract the special paramters from the uri provided in validation. These two styles are the Classic and Mapped pattern.

##Install

     go get github.com/influx6/reggy

  Then

    go install github.com/influx6/reggy


##API

- Matchable Interface{}
    This is the core interface of all Reggy matchers and they all must atleast meet the function requirements of this interface.

- ClassicMatcher struct{}
    This struct defines a single part of set of the given pattern i.e in /user/name/:id ,there will be a classicMatcher  that defines the rules for each of those pieces of the part. It is used when the first standard Reggy matching approach is used.

- FunctionMatcher struct{}
    This struct defines a single piece of a set of parts and is used in the creation of parts that provide custom functions to validate their validity.


- ClassicPattern struct{}
    This struct defines and contains all ClassicMatchers of the pieces of a given pattern and it contains the necessary logic for the uses of all of them for the validation of all pieces of the pattern

- MappedPattern struct{}
    This struct defines and contains all FunctionMatchers of the pieces of a given pattern and it contains the necessary logic for the uses of all of them for the validation of all pieces of the pattern

####ClassicMatchMux
These struct encapsulates the classic pattern and provides the pattern and parameter extracting features

- CreateClassic(pattern string)
    This creates the classic pattern matcher which takes the pattern with its mix of strings and regular expressions to define what stands as valid piece of the pattern

- ClassicMatchMux.Validate(pattern string,beStrictOnLength bool) (bool,map)
    This member function takes the pattern to match against and a bool value to indicate if should be strict on the length of the patterns to be matched.It returns a boolean and a map containing marked out pieces in the pattern, pieces are marked by the fact they have a ‘:[regular expression]’ in their definition

        pattern   = `/name/{id:[\d+]}`
        //or pattern   = `/name/:id` to match anything
        r := CreateClassic(pattern)

        state, param := r.Validate(`/name/12`, false)

        where:
            param contains the key points extracted according to the piece and its regular expression sections
            state is a boolean stating if it passed/true or failed/false


####FunctionalMatchMux
These struct encapsulates the mapped pattern and provides the pattern and parameter extracting features

- CreateFunctional(pattern string,map[string]interface{})
    This creates the mapped pattern matcher which takes the pattern and a map of named attributes and functions to decides what is valid or not

- FunctionalMatchMux.Validate(pattern string,beStrictOnLength bool) (bool,map)
        This member function takes the pattern to match against and a bool value to indicate if should be strict on the length of the patterns to be matched.It returns a boolean and a map containing marked out pieces in the pattern, pieces are marked by the fact they value functions to validate them in the pattern map

        pattern   = `/name/id`
        validators = MapFunc{
                "id": func(i interface{}) bool {
                rs := i.(string)
                if numbOnly.MatchString(rs) {
                    return true
                }
                return false
            },
        }

        r := CreateFunctional(pattern, validators)
        state, param := r.Validate(`/name/12`, false)

        where:
            param contains the key in the function map supplied
            state is a boolean stating if it passed/true or failed/false
